package app

import (
	"testing"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"github.com/stretchr/testify/assert"
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

	t.Run("basic unverified user can ", func(t *testing.T) {
		session := model.Session{
			UserId: th.BasicUnverified.Id,
		}
		assert.False(t, th.App.SessionHasPermissionToChannels(th.Context, session, allChannels, model.PermissionCreatePost))
	})
}
