package localcachelayer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store/storetest/mocks"
)

func TestBlocklistStoreChannelCache(t *testing.T) {
	channelId := "channel1"
	userId := "user1"
	fakeBlockList := model.ChannelBlockUserList{
			{
				ChannelId: channelId,
				BlockedId:    userId,
			},
	}

	t.Run("first call not cached, second cached and returning same data", func(t *testing.T) {
		mockStore := getMockStore(t)
		mockCacheProvider := getMockCacheProvider()
		cachedStore, err := NewLocalCacheLayer(mockStore, nil, nil, mockCacheProvider)
		require.NoError(t, err)

		blockList, err := cachedStore.Blocklist().ListChannelBlockUsers(channelId)
		require.NoError(t, err)
		assert.Equal(t, blockList, &fakeBlockList)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListChannelBlockUsers", 1)

		blockList, err = cachedStore.Blocklist().ListChannelBlockUsers(channelId)
		require.NoError(t, err)
		assert.Equal(t, blockList, &fakeBlockList)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListChannelBlockUsers", 1)
	})

	t.Run("first call not cached, invalidate cache, second call not cached", func(t *testing.T) {
		mockStore := getMockStore(t)
		mockCacheProvider := getMockCacheProvider()
		cachedStore, err := NewLocalCacheLayer(mockStore, nil, nil, mockCacheProvider)
		require.NoError(t, err)

		cachedStore.Blocklist().ListChannelBlockUsers(channelId)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListChannelBlockUsers", 1)
		
		cachedStore.Blocklist().InvalidateCacheForChannel(channelId)
		
		cachedStore.Blocklist().ListChannelBlockUsers(channelId)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListChannelBlockUsers", 2)
	})
}

func TestBlocklistStoreTeamCache(t *testing.T) {
	teamId := "team1"
	userId := "user1"
	fakeBlockList := model.TeamBlockUserList{
		{
			TeamId: teamId,
			BlockedId: userId,
		},
	}

	t.Run("first call not cached, second cached and returning same data", func(t *testing.T) {
		mockStore := getMockStore(t)
		mockCacheProvider := getMockCacheProvider()
		cachedStore, err := NewLocalCacheLayer(mockStore, nil, nil, mockCacheProvider)
		require.NoError(t, err)

		blockList, err := cachedStore.Blocklist().ListTeamBlockUsers(teamId)
		require.NoError(t, err)
		assert.Equal(t, blockList, &fakeBlockList)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListTeamBlockUsers", 1)

		blockList, err = cachedStore.Blocklist().ListTeamBlockUsers(teamId)
		require.NoError(t, err)
		assert.Equal(t, blockList, &fakeBlockList)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListTeamBlockUsers", 1)
	})

	t.Run("first call not cached, invalidate cache, second call not cached", func(t *testing.T) {
		mockStore := getMockStore(t)
		mockCacheProvider := getMockCacheProvider()
		cachedStore, err := NewLocalCacheLayer(mockStore, nil, nil, mockCacheProvider)
		require.NoError(t, err)

		cachedStore.Blocklist().ListTeamBlockUsers(teamId)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListTeamBlockUsers", 1)
		
		cachedStore.Blocklist().InvalidateCacheForTeam(teamId)
		
		cachedStore.Blocklist().ListTeamBlockUsers(teamId)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListTeamBlockUsers", 2)
	})
}

func TestBlocklistStoreUserCache(t *testing.T) {
	userId := "user1"
	blockedId := "user2"
	fakeBlockList := model.UserBlockUserList{
		{
			UserId:    userId,
			BlockedId: blockedId,
		},
	}

	t.Run("first call not cached, second cached and returning same data", func(t *testing.T) {
		mockStore := getMockStore(t)
		mockCacheProvider := getMockCacheProvider()
		cachedStore, err := NewLocalCacheLayer(mockStore, nil, nil, mockCacheProvider)
		require.NoError(t, err)

		blockList, err := cachedStore.Blocklist().ListUserBlockUsers(userId)
		require.NoError(t, err)
		assert.Equal(t, blockList, &fakeBlockList)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListUserBlockUsers", 1)

		blockList, err = cachedStore.Blocklist().ListUserBlockUsers(userId)
		require.NoError(t, err)
		assert.Equal(t, blockList, &fakeBlockList)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListUserBlockUsers", 1)
	})

	t.Run("first call not cached, invalidate cache, second call not cached", func(t *testing.T) {
		mockStore := getMockStore(t)
		mockCacheProvider := getMockCacheProvider()
		cachedStore, err := NewLocalCacheLayer(mockStore, nil, nil, mockCacheProvider)
		require.NoError(t, err)

		cachedStore.Blocklist().ListUserBlockUsers(userId)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListUserBlockUsers", 1)
		
		cachedStore.Blocklist().InvalidateCacheForUser(userId)
		
		cachedStore.Blocklist().ListUserBlockUsers(userId)
		mockStore.Blocklist().(*mocks.BlocklistStore).AssertNumberOfCalls(t, "ListUserBlockUsers", 2)
	})
}