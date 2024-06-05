package storetest

import (
	"testing"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannelBlockUserStore(t *testing.T, rctx request.CTX, ss store.Store) {
	t.Run("SaveChannelBlockUser", func(t *testing.T) { testSaveChannelBlockUser(t, rctx, ss) })
}

func testSaveChannelBlockUser(t *testing.T, _ request.CTX, ss store.Store) {
	channelBlockUser := model.ChannelBlockUser{
		ChannelId: "000000",
		BlockedId: "abcdefg",
		CreateBy:  "abcdfg",
	}

	var err error
	var newBlockUser, getBlockUser *model.ChannelBlockUser
	newBlockUser, err = ss.Blocklist().SaveChannelBlockUser(&channelBlockUser)
	require.NoError(t, err)
	require.Greater(t, newBlockUser.CreateAt, model.GetMillis()-200)
	assert.Equal(t, channelBlockUser.BlockedId, newBlockUser.BlockedId)
	assert.Equal(t, channelBlockUser.ChannelId, newBlockUser.ChannelId)
	assert.Equal(t, channelBlockUser.CreateBy, newBlockUser.CreateBy)

	getBlockUser, err = ss.Blocklist().GetChannelBlockUser(channelBlockUser.ChannelId, channelBlockUser.BlockedId)
	require.NoError(t, err)
	assert.Equal(t, channelBlockUser.BlockedId, getBlockUser.BlockedId)
	assert.Equal(t, channelBlockUser.ChannelId, getBlockUser.ChannelId)
	assert.Equal(t, channelBlockUser.CreateBy, getBlockUser.CreateBy)
}
