package slashcommands

import (
	"testing"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"github.com/stretchr/testify/assert"
)


func TestAnnouncProviderCommand(t *testing.T) {
	th := setup(t).initBasic()
	defer th.tearDown()
	newUser := th.createUser()
	th.linkUserToTeam(newUser, th.BasicTeam)
	th.addUserToChannel(newUser, th.BasicChannel)
	
	ok := th.App.HasPermissionToChannel(th.Context, newUser.Id, th.BasicChannel.Id, model.PermissionCreatePost)
	assert.True(t, ok)

	ap := AnnouncProvider{}

	args := &model.CommandArgs{
		T:         func(s string, args ...any) string { return s },
		ChannelId: th.BasicChannel.Id,
		UserId:    th.BasicUser.Id,
	}

	th.removePermissionFromRole(model.PermissionManageChannelRoles.Id, model.ChannelAdminRoleId)

	resp := ap.DoCommand(th.App, th.Context, args, "").Text
	assert.Equal(t, "api.command_channel_announc.permission.app_error", resp)

	th.addPermissionToRole(model.PermissionManageChannelRoles.Id, model.ChannelAdminRoleId)
	// empty message
	resp = ap.DoCommand(th.App, th.Context, args, "").Text
	assert.Equal(t, "api.command_channel_announc.message.app_error", resp)

	// random message
	resp = ap.DoCommand(th.App, th.Context, args, "random").Text
	assert.Equal(t, "api.command_channel_announc.param.param_error", resp)

	resp = ap.DoCommand(th.App, th.Context, args, "yes").Text
	assert.Equal(t, "", resp)
	ok = th.App.HasPermissionToChannel(th.Context, newUser.Id, th.BasicChannel.Id, model.PermissionCreatePost)
	assert.False(t, ok)

	resp = ap.DoCommand(th.App, th.Context, args, "no").Text
	assert.Equal(t, "", resp)
	ok = th.App.HasPermissionToChannel(th.Context, newUser.Id, th.BasicChannel.Id, model.PermissionCreatePost)
	assert.True(t, ok)
}