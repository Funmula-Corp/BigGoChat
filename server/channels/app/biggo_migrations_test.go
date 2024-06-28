package app

import (
	"testing"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnverifiedSessionHasPermissionToChannel(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	ch1 := th.CreateChannel(th.Context, th.BasicTeam)
	ch2 := th.CreatePrivateChannel(th.Context, th.BasicTeam)
	th.App.AddUserToChannel(th.Context, th.BasicUnverified, ch1, false)
	th.App.AddUserToChannel(th.Context, th.BasicUnverified, ch2, false)
	th.App.AddUserToChannel(th.Context, th.BasicUnverified, th.BasicChannel, false)
	allChannels := []string{th.BasicChannel.Id, ch1.Id, ch2.Id}

	t.Run("basic unverified user can read basic channels", func(t *testing.T) {
		session := model.Session{
			UserId: th.BasicUnverified.Id,
		}
		assert.True(t, th.App.SessionHasPermissionToChannels(th.Context, session, allChannels, model.PermissionReadChannel))
	})

	t.Run("basic unverified user can not post", func(t *testing.T) {
		session := model.Session{
			UserId: th.BasicUnverified.Id,
		}
		assert.False(t, th.App.SessionHasPermissionToChannels(th.Context, session, allChannels, model.PermissionCreatePost))
	})

	t.Run("system verified user can create team", func(t *testing.T){
		model.IsInRole(th.BasicUser.Roles, SystemVerifiedRoleId)
		role, err := th.App.GetRoleByName(th.Context.Context(), model.SystemVerifiedRoleId)
		require.Nil(t, err)
		roleMap := make(map[string]bool)
		for _, permission := range(role.Permissions){
			roleMap[permission] = true
		}
		require.True(t, roleMap[model.PermissionCreateTeam.Id])

		session := model.Session{
			UserId: th.BasicUser.Id,
			Roles: th.BasicUser.Roles,
		}
		require.True(t, th.App.SessionHasPermissionTo(session, model.PermissionCreateTeam))
	})

	t.Run("basic unverified user can not create team", func(t *testing.T) {
		session := model.Session{
			UserId: th.BasicUnverified.Id,
			Roles: th.BasicUnverified.Roles,
		}
		assert.False(t, th.App.SessionHasPermissionTo(session, model.PermissionCreateTeam))
	})
}

func TestMigratedChannelPrivacyConvert(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	t.Run("teamadmin convert private to public", func(t *testing.T){
		cm, err := th.App.GetChannelMember(th.Context, th.BasicChannel.Id, th.BasicUser.Id)
		require.Nil(t, err)
		require.True(t, cm.SchemeAdmin)
		session := model.Session {
			UserId: th.BasicUser.Id,
		}
		require.False(t, th.App.SessionHasPermissionToChannels(th.Context, session, []string{th.BasicChannel.Id}, model.PermissionConvertPrivateChannelToPublic))
		require.True(t, th.App.SessionHasPermissionToChannels(th.Context, session, []string{th.BasicChannel.Id}, model.PermissionConvertPublicChannelToPrivate))

		session = model.Session {
			UserId: th.BasicUser2.Id,
		}
		require.False(t, th.App.SessionHasPermissionToChannels(th.Context, session, []string{th.BasicChannel.Id}, model.PermissionConvertPrivateChannelToPublic))
		require.False(t, th.App.SessionHasPermissionToChannels(th.Context, session, []string{th.BasicChannel.Id}, model.PermissionConvertPublicChannelToPrivate))
		id := model.NewId()
		channel := &model.Channel {
			DisplayName: "dn_" + id,
			Name:        "name_" + id,
			Type:        model.ChannelTypePrivate,
			TeamId:      th.BasicTeam.Id,
			CreatorId:   th.BasicUser2.Id,
		}
		channel, err = th.App.CreateChannel(th.Context, channel, true)
		require.Nil(t, err)
		require.False(t, th.App.SessionHasPermissionToChannels(th.Context, session, []string{channel.Id}, model.PermissionConvertPrivateChannelToPublic))
		require.True(t, th.App.SessionHasPermissionToChannels(th.Context, session, []string{channel.Id}, model.PermissionConvertPublicChannelToPrivate))

		// BasicUser2 can convert BasicChannel to private after promotion to a team moderator.
		_, err = th.App.UpdateTeamMemberSchemeRoles(th.Context, th.BasicTeam.Id, th.BasicUser2.Id, false, true, true, true, false)
		require.Nil(t, err)
		sessions, err := th.App.GetSessions(th.Context, th.BasicUser2.Id)
		require.Nil(t, err)
		for _, s := range(sessions){
			require.False(t, th.App.SessionHasPermissionToChannels(th.Context, *s, []string{th.BasicChannel.Id}, model.PermissionConvertPrivateChannelToPublic))
			require.True(t, th.App.SessionHasPermissionToChannels(th.Context, *s, []string{th.BasicChannel.Id}, model.PermissionConvertPublicChannelToPrivate))
		}
	})
}
