package slashcommands

import (
	"strings"
	"testing"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"github.com/stretchr/testify/assert"
)

func TestMuteuserProviderCommand(t *testing.T) {
	th := setup(t).initBasic()
	defer th.tearDown()
	user1 := th.createUser()
	user2 := th.createUser()
	users := []*model.User{user1, user2}

	for _, user := range users {
		th.linkUserToTeam(user, th.BasicTeam)
		th.addUserToChannel(user, th.BasicChannel)
		ok := th.App.HasPermissionToChannel(th.Context, user.Id, th.BasicChannel.Id, model.PermissionCreatePost)
		assert.True(t, ok)
	}

	mp := MuteuserProvider{}
	ump := UnMuteuserProvider{}

	args := &model.CommandArgs{
		T:         func(s string, args ...any) string { return s },
		ChannelId: th.BasicChannel.Id,
		UserId:    th.BasicUser.Id,
	}

	th.removePermissionFromRole(model.PermissionManageChannelRoles.Id, model.ChannelAdminRoleId)

	resp := mp.DoCommand(th.App, th.Context, args, "").Text
	assert.Equal(t, "api.command_channel_muteuser.permission.app_error", resp)

	th.addPermissionToRole(model.PermissionManageChannelRoles.Id, model.ChannelAdminRoleId)
	// empty message
	resp = mp.DoCommand(th.App, th.Context, args, "").Text
	assert.Equal(t, "api.command_channel_muteuser.param.length_error", resp)

	// random message
	resp = mp.DoCommand(th.App, th.Context, args, "random").Text
	assert.Equal(t, "api.command_channel_muteuser.param.length_error", resp)

	resp = mp.DoCommand(th.App, th.Context, args, "@NonExistsUser").Text
	assert.Equal(t, "api.command_channel_muteuser.param.length_error", resp)

	// mute single user test
	for _, user := range users {
		resp = mp.DoCommand(th.App, th.Context, args, "@"+user.Username).Text
		assert.Equal(t, "", resp)
		ok := th.App.HasPermissionToChannel(th.Context, user.Id, th.BasicChannel.Id, model.PermissionCreatePost)
		assert.False(t, ok)
	}

	// unmute single user test
	for _, user := range users {
		resp = ump.DoCommand(th.App, th.Context, args, "@"+user.Username).Text
		assert.Equal(t, "", resp)
		ok := th.App.HasPermissionToChannel(th.Context, user.Id, th.BasicChannel.Id, model.PermissionCreatePost)
		assert.True(t, ok)
	}

	// mute multiple user
	message := strings.Join([]string{"@"+user1.Username, "@"+user2.Username}, " ")
	resp = mp.DoCommand(th.App, th.Context, args, message).Text
	assert.Equal(t, "", resp)
	for _, user := range users {
		ok := th.App.HasPermissionToChannel(th.Context, user.Id, th.BasicChannel.Id, model.PermissionCreatePost)
		assert.False(t, ok)
	}

	// unmute multiple user
	resp = ump.DoCommand(th.App, th.Context, args, message).Text
	assert.Equal(t, "", resp)
	for _, user := range users {
		ok := th.App.HasPermissionToChannel(th.Context, user.Id, th.BasicChannel.Id, model.PermissionCreatePost)
		assert.True(t, ok)
	}
}
