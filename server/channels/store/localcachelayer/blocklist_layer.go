package localcachelayer

import (
	"bytes"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store"
)

type LocalCacheBlocklistStore struct {
	store.BlocklistStore
	rootStore *LocalCacheStore
}

func (s *LocalCacheBlocklistStore) handleClusterInvalidateChannel(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.channelBlocklistCache.Purge()
	} else {
		s.rootStore.channelBlocklistCache.Remove(string(msg.Data))
	}
}

func (s *LocalCacheBlocklistStore) handleClusterInvalidateUser(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.userBlocklistCache.Purge()
	} else {
		s.rootStore.userBlocklistCache.Remove(string(msg.Data))
	}
}

func (s *LocalCacheBlocklistStore) handleClusterInvalidateTeam(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.teamBlocklistCache.Purge()
	} else {
		s.rootStore.teamBlocklistCache.Remove(string(msg.Data))
	}
}

func (s LocalCacheBlocklistStore) InvalidateCacheForChannel(channelId string) {
	s.rootStore.doInvalidateCacheCluster(s.rootStore.channelBlocklistCache, channelId, nil)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelBlocklistCache.Name())
	}
}

func (s LocalCacheBlocklistStore) InvalidateCacheForUser(userId string) {
	s.rootStore.doInvalidateCacheCluster(s.rootStore.userBlocklistCache, userId, nil)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.userBlocklistCache.Name())
	}
}

func (s LocalCacheBlocklistStore) InvalidateCacheForTeam(teamId string) {
	s.rootStore.doInvalidateCacheCluster(s.rootStore.teamBlocklistCache, teamId, nil)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.teamBlocklistCache.Name())
	}
}


func (s LocalCacheBlocklistStore) ListChannelBlockUsers(channelId string) (*model.ChannelBlockUserList, error) {
	var blockUserList model.ChannelBlockUserList
	if err := s.rootStore.doStandardReadCache(s.rootStore.channelBlocklistCache, channelId, &blockUserList); err == nil {
		return &blockUserList, nil
	}
	blockUserListPtr, err := s.BlocklistStore.ListChannelBlockUsers(channelId)
	if err != nil {
		return nil, err
	}
	s.rootStore.doStandardAddToCache(s.rootStore.channelBlocklistCache, channelId, *blockUserListPtr)
	return blockUserListPtr, nil
}

func (s LocalCacheBlocklistStore) ListTeamBlockUsers(channelId string) (*model.TeamBlockUserList, error) {
	var blockUserList model.TeamBlockUserList
	if err := s.rootStore.doStandardReadCache(s.rootStore.teamBlocklistCache, channelId, &blockUserList); err == nil {
		return &blockUserList, nil
	}
	blockUserListPtr, err := s.BlocklistStore.ListTeamBlockUsers(channelId)
	if err != nil {
		return nil, err
	}
	s.rootStore.doStandardAddToCache(s.rootStore.teamBlocklistCache, channelId, *blockUserListPtr)
	return blockUserListPtr, nil
}

func (s LocalCacheBlocklistStore) ListUserBlockUsers(userId string) (*model.UserBlockUserList, error) {
	var blockUserList model.UserBlockUserList
	if err := s.rootStore.doStandardReadCache(s.rootStore.userBlocklistCache, userId, &blockUserList); err == nil {
		return &blockUserList, nil
	}
	blockUserListPtr, err := s.BlocklistStore.ListUserBlockUsers(userId)
	if err != nil {
		return nil, err
	}
	s.rootStore.doStandardAddToCache(s.rootStore.userBlocklistCache, userId, *blockUserListPtr)
	return blockUserListPtr, nil
}

func (s LocalCacheBlocklistStore) DeleteChannelBlockUser(channelId string, userId string) error {
	err := s.BlocklistStore.DeleteChannelBlockUser(channelId, userId)
	if err != nil {
		return err
	}
	s.InvalidateCacheForChannel(channelId)
	return nil
}

func (s LocalCacheBlocklistStore) DeleteTeamBlockUser(channelId string, userId string) error {
	err := s.BlocklistStore.DeleteTeamBlockUser(channelId, userId)
	if err != nil {
		return err
	}
	s.InvalidateCacheForTeam(channelId)
	return nil
}

func (s LocalCacheBlocklistStore) DeleteUserBlockUser(userId, blockedId string, userIsVerified, blockedIsVerified bool) error {
	err := s.BlocklistStore.DeleteUserBlockUser(userId, blockedId, userIsVerified, blockedIsVerified)
	if err != nil {
		return err
	}
	s.InvalidateCacheForUser(userId)
	return nil
}

func (s LocalCacheBlocklistStore) SaveChannelBlockUser(blockUser *model.ChannelBlockUser) (*model.ChannelBlockUser, error) {
	blockUser, err := s.BlocklistStore.SaveChannelBlockUser(blockUser)
	if err != nil {
		return nil, err
	}
	s.InvalidateCacheForChannel(blockUser.ChannelId)
	return blockUser, nil
}

func (s LocalCacheBlocklistStore) SaveTeamBlockUser(blockUser *model.TeamBlockUser) (*model.TeamBlockUser, error) {
	blockUser, err := s.BlocklistStore.SaveTeamBlockUser(blockUser)
	if err != nil {
		return nil, err
	}
	s.InvalidateCacheForTeam(blockUser.TeamId)
	return blockUser, nil
}

func (s LocalCacheBlocklistStore) SaveUserBlockUser(userBlockUser *model.UserBlockUser) (*model.UserBlockUser, error) {
	userBlockUser, err := s.BlocklistStore.SaveUserBlockUser(userBlockUser)
	if err != nil {
		return nil, err
	}
	s.InvalidateCacheForUser(userBlockUser.UserId)
	return userBlockUser, nil
}
