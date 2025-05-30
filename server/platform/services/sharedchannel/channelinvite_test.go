// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/plugin/plugintest/mock"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store/storetest/mocks"
)

var (
	mockTypeChannel    = mock.AnythingOfType("*model.Channel")
	mockTypeString     = mock.AnythingOfType("string")
	mockTypeReqContext = mock.AnythingOfType("*request.Context")
	mockTypeContext    = mock.MatchedBy(func(ctx context.Context) bool { return true })
)

func TestOnReceiveChannelInvite(t *testing.T) {
	t.Run("when msg payload is empty, it does nothing", func(t *testing.T) {
		mockServer := &MockServerIface{}
		logger := mlog.CreateConsoleTestLogger(t)
		mockServer.On("Log").Return(logger)
		mockApp := &MockAppIface{}
		scs := &Service{
			server: mockServer,
			app:    mockApp,
		}

		mockStore := &mocks.Store{}
		mockServer = scs.server.(*MockServerIface)
		mockServer.On("GetStore").Return(mockStore)

		remoteCluster := &model.RemoteCluster{}
		msg := model.RemoteClusterMsg{}

		err := scs.onReceiveChannelInvite(msg, remoteCluster, nil)
		require.NoError(t, err)
		mockStore.AssertNotCalled(t, "Channel")
	})

	t.Run("when invitation prescribes a readonly channel, it does create a readonly channel", func(t *testing.T) {
		mockServer := &MockServerIface{}
		logger := mlog.CreateConsoleTestLogger(t)
		mockServer.On("Log").Return(logger)
		mockApp := &MockAppIface{}
		scs := &Service{
			server: mockServer,
			app:    mockApp,
		}

		mockStore := &mocks.Store{}
		remoteCluster := &model.RemoteCluster{Name: "test"}
		invitation := channelInviteMsg{
			ChannelId: model.NewId(),
			TeamId:    model.NewId(),
			ReadOnly:  true,
			Type:      model.ChannelTypeOpen,
		}
		payload, err := json.Marshal(invitation)
		require.NoError(t, err)

		msg := model.RemoteClusterMsg{
			Payload: payload,
		}
		mockChannelStore := mocks.ChannelStore{}
		mockSharedChannelStore := mocks.SharedChannelStore{}
		channel := &model.Channel{
			Id:     invitation.ChannelId,
			TeamId: invitation.TeamId,
			Type:   invitation.Type,
		}

		mockChannelStore.On("Get", invitation.ChannelId, true).Return(nil, &store.ErrNotFound{})
		mockSharedChannelStore.On("Save", mock.Anything).Return(nil, nil)
		mockSharedChannelStore.On("SaveRemote", mock.Anything).Return(nil, nil)
		mockStore.On("Channel").Return(&mockChannelStore)
		mockStore.On("SharedChannel").Return(&mockSharedChannelStore)

		mockServer.On("GetStore").Return(mockStore)
		createPostPermission := model.ChannelModeratedPermissionsMap[model.PermissionCreatePost.Id]
		createReactionPermission := model.ChannelModeratedPermissionsMap[model.PermissionAddReaction.Id]
		updateMap := model.ChannelModeratedRolesPatch{
			Guests:  model.NewBool(false),
			Members: model.NewBool(false),
		}

		mockApp.On("CreateChannelWithUser", mockTypeReqContext, mockTypeChannel, mockTypeString).Return(channel, nil)

		readonlyChannelModerations := []*model.ChannelModerationPatch{
			{
				Name:  &createPostPermission,
				Roles: &updateMap,
			},
			{
				Name:  &createReactionPermission,
				Roles: &updateMap,
			},
		}
		mockApp.On("PatchChannelModerationsForChannel", mock.Anything, channel, readonlyChannelModerations).Return(nil, nil).Maybe()
		defer mockApp.AssertExpectations(t)

		err = scs.onReceiveChannelInvite(msg, remoteCluster, nil)
		require.NoError(t, err)
	})

	t.Run("when invitation prescribes a readonly channel and readonly update fails, it returns an error", func(t *testing.T) {
		mockServer := &MockServerIface{}
		logger := mlog.CreateConsoleTestLogger(t)
		mockServer.On("Log").Return(logger)
		mockApp := &MockAppIface{}
		scs := &Service{
			server: mockServer,
			app:    mockApp,
		}

		mockStore := &mocks.Store{}
		remoteCluster := &model.RemoteCluster{Name: "test2"}
		invitation := channelInviteMsg{
			ChannelId: model.NewId(),
			TeamId:    model.NewId(),
			ReadOnly:  true,
			Type:      "0",
		}
		payload, err := json.Marshal(invitation)
		require.NoError(t, err)

		msg := model.RemoteClusterMsg{
			Payload: payload,
		}
		mockChannelStore := mocks.ChannelStore{}
		channel := &model.Channel{
			Id: invitation.ChannelId,
		}

		mockChannelStore.On("Get", invitation.ChannelId, true).Return(nil, &store.ErrNotFound{})
		mockStore.On("Channel").Return(&mockChannelStore)

		mockServer = scs.server.(*MockServerIface)
		mockServer.On("GetStore").Return(mockStore)
		appErr := model.NewAppError("foo", "bar", nil, "boom", http.StatusBadRequest)

		mockApp.On("CreateChannelWithUser", mockTypeReqContext, mockTypeChannel, mockTypeString).Return(channel, nil)
		mockApp.On("PatchChannelModerationsForChannel", mock.Anything, channel, mock.Anything).Return(nil, appErr)
		defer mockApp.AssertExpectations(t)

		err = scs.onReceiveChannelInvite(msg, remoteCluster, nil)
		require.Error(t, err)
		assert.Equal(t, fmt.Sprintf("cannot make channel readonly `%s`: foo: bar, boom", invitation.ChannelId), err.Error())
	})

	t.Run("DM channels", func(t *testing.T) {
		var testRemoteID = model.NewId()
		testCases := []struct {
			desc          string
			user1         *model.User
			user2         *model.User
			canSee        bool
			expectSuccess bool
		}{
			{"valid users", &model.User{Id: model.NewId(), RemoteId: &testRemoteID}, &model.User{Id: model.NewId()}, true, true},
			{"swapped users", &model.User{Id: model.NewId()}, &model.User{Id: model.NewId(), RemoteId: &testRemoteID}, true, true},
			{"two remotes", &model.User{Id: model.NewId(), RemoteId: &testRemoteID}, &model.User{Id: model.NewId(), RemoteId: &testRemoteID}, true, false},
			{"two locals", &model.User{Id: model.NewId()}, &model.User{Id: model.NewId()}, true, false},
			{"can't see", &model.User{Id: model.NewId(), RemoteId: &testRemoteID}, &model.User{Id: model.NewId()}, false, false},
			{"invalid remoteid", &model.User{Id: model.NewId(), RemoteId: model.NewString("bogus")}, &model.User{Id: model.NewId()}, true, false},
		}

		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				mockServer := &MockServerIface{}
				logger := mlog.CreateConsoleTestLogger(t)
				mockServer.On("Log").Return(logger)
				mockApp := &MockAppIface{}
				scs := &Service{
					server: mockServer,
					app:    mockApp,
				}

				mockStore := &mocks.Store{}
				remoteCluster := &model.RemoteCluster{Name: "test3", CreatorId: model.NewId(), RemoteId: testRemoteID}
				invitation := channelInviteMsg{
					ChannelId:            model.NewId(),
					TeamId:               model.NewId(),
					ReadOnly:             false,
					Type:                 model.ChannelTypeDirect,
					DirectParticipantIDs: []string{tc.user1.Id, tc.user2.Id},
				}
				payload, err := json.Marshal(invitation)
				require.NoError(t, err)

				msg := model.RemoteClusterMsg{
					Payload: payload,
				}
				mockChannelStore := mocks.ChannelStore{}
				mockSharedChannelStore := mocks.SharedChannelStore{}
				channel := &model.Channel{
					Id: invitation.ChannelId,
				}

				mockUserStore := mocks.UserStore{}
				mockUserStore.On("Get", mockTypeContext, tc.user1.Id).
					Return(tc.user1, nil)
				mockUserStore.On("Get", mockTypeContext, tc.user2.Id).
					Return(tc.user2, nil)

				mockChannelStore.On("Get", invitation.ChannelId, true).Return(nil, errors.New("boom"))
				mockChannelStore.On("GetByName", "", mockTypeString, true).Return(nil, &store.ErrNotFound{})

				mockSharedChannelStore.On("Save", mock.Anything).Return(nil, nil)
				mockSharedChannelStore.On("SaveRemote", mock.Anything).Return(nil, nil)
				mockStore.On("Channel").Return(&mockChannelStore)
				mockStore.On("SharedChannel").Return(&mockSharedChannelStore)
				mockStore.On("User").Return(&mockUserStore)

				mockServer = scs.server.(*MockServerIface)
				mockServer.On("GetStore").Return(mockStore)

				mockApp.On("GetOrCreateDirectChannel", mockTypeReqContext, mockTypeString, mockTypeString, mock.AnythingOfType("model.ChannelOption")).
					Return(channel, nil).Maybe()
				mockApp.On("UserCanSeeOtherUser", mockTypeReqContext, mockTypeString, mockTypeString).Return(tc.canSee, nil).Maybe()

				defer mockApp.AssertExpectations(t)

				err = scs.onReceiveChannelInvite(msg, remoteCluster, nil)
				require.Equal(t, tc.expectSuccess, err == nil)
			})
		}
	})
}
