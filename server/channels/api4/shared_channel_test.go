// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func setupForSharedChannels(tb testing.TB) *TestHelper {
	return SetupConfig(tb, func(cfg *model.Config) {
		*cfg.ExperimentalSettings.EnableRemoteClusterService = true
		*cfg.ExperimentalSettings.EnableSharedChannels = true
	})
}

func TestGetAllSharedChannels(t *testing.T) {
	th := setupForSharedChannels(t).InitBasic()
	defer th.TearDown()

	const pages = 3
	const pageSize = 7

	savedIds := make([]string, 0, pages*pageSize)

	// make some shared channels
	for i := 0; i < pages*pageSize; i++ {
		channel := th.CreateChannelWithClientAndTeam(th.Client, model.ChannelTypeOpen, th.BasicTeam.Id)
		sc := &model.SharedChannel{
			ChannelId: channel.Id,
			TeamId:    channel.TeamId,
			Home:      randomBool(),
			ShareName: fmt.Sprintf("test_share_%d", i),
			CreatorId: th.BasicChannel.CreatorId,
			RemoteId:  model.NewId(),
		}

		_, err := th.App.ShareChannel(th.Context, sc)
		require.NoError(t, err)
		savedIds = append(savedIds, channel.Id)
	}
	sort.Strings(savedIds)

	t.Run("get shared channels paginated", func(t *testing.T) {
		channelIds := make([]string, 0, 21)
		for i := 0; i < pages; i++ {
			channels, _, err := th.Client.GetAllSharedChannels(context.Background(), th.BasicTeam.Id, i, pageSize)
			require.NoError(t, err)
			channelIds = append(channelIds, getIds(channels)...)
		}
		sort.Strings(channelIds)

		// ids lists should now match
		assert.Equal(t, savedIds, channelIds, "id lists should match")
	})

	t.Run("get shared channels for invalid team", func(t *testing.T) {
		_, _, err := th.Client.GetAllSharedChannels(context.Background(), model.NewId(), 0, 100)
		require.Error(t, err)
	})

	t.Run("get shared channels, user not member of team", func(t *testing.T) {
		team := &model.Team{
			DisplayName: "tteam",
			Name:        GenerateTestTeamName(),
			Type:        model.TeamOpen,
		}
		team, _, err := th.SystemAdminClient.CreateTeam(context.Background(), team)
		require.NoError(t, err)

		_, _, err = th.Client.GetAllSharedChannels(context.Background(), team.Id, 0, 100)
		require.Error(t, err)
	})
}

func getIds(channels []*model.SharedChannel) []string {
	ids := make([]string, 0, len(channels))
	for _, c := range channels {
		ids = append(ids, c.ChannelId)
	}
	return ids
}

func randomBool() bool {
	return rnd.Intn(2) != 0
}

func TestGetRemoteClusterById(t *testing.T) {
	th := setupForSharedChannels(t).InitBasic()
	defer th.TearDown()

	// for this test we need a user that belongs to a channel that
	// is shared with the requested remote id.

	// create a remote cluster
	rc := &model.RemoteCluster{
		RemoteId:  model.NewId(),
		Name:      "Test1",
		SiteURL:   model.NewId(),
		CreatorId: model.NewId(),
	}
	rc, appErr := th.App.AddRemoteCluster(rc)
	require.Nil(t, appErr)

	// create a shared channel
	sc := &model.SharedChannel{
		ChannelId: th.BasicChannel.Id,
		TeamId:    th.BasicChannel.TeamId,
		Home:      false,
		ShareName: "test_share",
		CreatorId: th.BasicChannel.CreatorId,
		RemoteId:  rc.RemoteId,
	}
	sc, err := th.App.ShareChannel(th.Context, sc)
	require.NoError(t, err)

	// create a shared channel remote to connect them
	scr := &model.SharedChannelRemote{
		Id:                model.NewId(),
		ChannelId:         sc.ChannelId,
		CreatorId:         sc.CreatorId,
		IsInviteAccepted:  true,
		IsInviteConfirmed: true,
		RemoteId:          sc.RemoteId,
	}
	_, err = th.App.SaveSharedChannelRemote(scr)
	require.NoError(t, err)

	t.Run("valid remote, user is member", func(t *testing.T) {
		rcInfo, _, err := th.Client.GetRemoteClusterInfo(context.Background(), rc.RemoteId)
		require.NoError(t, err)
		assert.Equal(t, rc.Name, rcInfo.Name)
	})

	t.Run("invalid remote", func(t *testing.T) {
		_, resp, err := th.Client.GetRemoteClusterInfo(context.Background(), model.NewId())
		require.Error(t, err)
		CheckNotFoundStatus(t, resp)
	})
}

func TestCreateDirectChannelWithRemoteUser(t *testing.T) {
	t.Run("creates a local DM channel that is shared", func(t *testing.T) {
		th := setupForSharedChannels(t).InitBasic()
		defer th.TearDown()
		client := th.Client
		defer client.Logout(context.Background())

		localUser := th.BasicUser
		remoteUser := th.CreateUser()
		remoteUser.RemoteId = model.NewString(model.NewId())
		remoteUser, appErr := th.App.UpdateUser(th.Context, remoteUser, false)
		require.Nil(t, appErr)

		dm, _, err := client.CreateDirectChannel(context.Background(), localUser.Id, remoteUser.Id)
		require.NoError(t, err)

		channelName := model.GetDMNameFromIds(localUser.Id, remoteUser.Id)
		require.Equal(t, channelName, dm.Name, "dm name didn't match")
		assert.True(t, dm.IsShared())
	})

	t.Run("sends a shared channel invitation to the remote", func(t *testing.T) {
		th := setupForSharedChannels(t).InitBasic()
		defer th.TearDown()
		client := th.Client
		defer client.Logout(context.Background())

		localUser := th.BasicUser
		remoteUser := th.CreateUser()

		rc := &model.RemoteCluster{
			Name:      "test",
			Token:     model.NewId(),
			CreatorId: localUser.Id,
		}
		rc, appErr := th.App.AddRemoteCluster(rc)
		require.Nil(t, appErr)

		remoteUser.RemoteId = model.NewString(rc.RemoteId)
		remoteUser, appErr = th.App.UpdateUser(th.Context, remoteUser, false)
		require.Nil(t, appErr)

		dm, _, err := client.CreateDirectChannel(context.Background(), localUser.Id, remoteUser.Id)
		require.NoError(t, err)

		channelName := model.GetDMNameFromIds(localUser.Id, remoteUser.Id)
		require.Equal(t, channelName, dm.Name, "dm name didn't match")
		require.True(t, dm.IsShared())
	})

	t.Run("does not send a shared channel invitation to the remote when creator is remote", func(t *testing.T) {
		th := setupForSharedChannels(t).InitBasic()
		defer th.TearDown()
		client := th.Client
		defer client.Logout(context.Background())

		localUser := th.BasicUser
		remoteUser := th.CreateUser()

		rc := &model.RemoteCluster{
			Name:      "test",
			Token:     model.NewId(),
			CreatorId: localUser.Id,
		}
		rc, appErr := th.App.AddRemoteCluster(rc)
		require.Nil(t, appErr)

		remoteUser.RemoteId = model.NewString(rc.RemoteId)
		remoteUser, appErr = th.App.UpdateUser(th.Context, remoteUser, false)
		require.Nil(t, appErr)

		dm, _, err := client.CreateDirectChannel(context.Background(), remoteUser.Id, localUser.Id)
		require.NoError(t, err)

		channelName := model.GetDMNameFromIds(localUser.Id, remoteUser.Id)
		require.Equal(t, channelName, dm.Name, "dm name didn't match")
		require.True(t, dm.IsShared())
	})
}
