// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package localcachelayer

import (
	"bytes"
	"errors"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store"
	"git.biggo.com/Funmula/BigGoChat/server/v8/platform/services/cache"
)

type LocalCacheChannelStore struct {
	store.ChannelStore
	rootStore *LocalCacheStore
}

func (s *LocalCacheChannelStore) handleClusterInvalidateChannelMemberCounts(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.channelMemberCountsCache.Purge()
	} else {
		s.rootStore.channelMemberCountsCache.Remove(string(msg.Data))
	}
}

func (s *LocalCacheChannelStore) handleClusterInvalidateChannelPinnedPostCount(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.channelPinnedPostCountsCache.Purge()
	} else {
		s.rootStore.channelPinnedPostCountsCache.Remove(string(msg.Data))
	}
}

func (s *LocalCacheChannelStore) handleClusterInvalidateChannelGuestCounts(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.channelGuestCountCache.Purge()
	} else {
		s.rootStore.channelGuestCountCache.Remove(string(msg.Data))
	}
}

func (s *LocalCacheChannelStore) handleClusterInvalidateChannelById(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.channelByIdCache.Purge()
	} else {
		s.rootStore.channelByIdCache.Remove(string(msg.Data))
	}
}

func (s *LocalCacheChannelStore) handleClusterInvalidateChannelForUser(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.channelMembersForUserCache.Purge()
	} else {
		s.rootStore.channelMembersForUserCache.Remove(string(msg.Data))
	}
}

func (s *LocalCacheChannelStore) handleClusterInvalidateChannelMembersNotifyProps(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.channelMembersNotifyPropsCache.Purge()
	} else {
		s.rootStore.channelMembersNotifyPropsCache.Remove(string(msg.Data))
	}
}

func (s *LocalCacheChannelStore) handleClusterInvalidateChannelByName(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.channelByNameCache.Purge()
	} else {
		s.rootStore.channelByNameCache.Remove(string(msg.Data))
	}
}

func (s LocalCacheChannelStore) ClearMembersForUserCache() {
	s.rootStore.doClearCacheCluster(s.rootStore.channelMembersForUserCache)
}

func (s LocalCacheChannelStore) ClearCaches() {
	s.rootStore.doClearCacheCluster(s.rootStore.channelMemberCountsCache)
	s.rootStore.doClearCacheCluster(s.rootStore.channelPinnedPostCountsCache)
	s.rootStore.doClearCacheCluster(s.rootStore.channelGuestCountCache)
	s.rootStore.doClearCacheCluster(s.rootStore.channelByIdCache)
	s.rootStore.doClearCacheCluster(s.rootStore.channelMembersForUserCache)
	s.rootStore.doClearCacheCluster(s.rootStore.channelMembersNotifyPropsCache)
	s.rootStore.doClearCacheCluster(s.rootStore.channelByNameCache)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelMemberCountsCache.Name())
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelPinnedPostCountsCache.Name())
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelGuestCountCache.Name())
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelByIdCache.Name())
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelMembersForUserCache.Name())
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelMembersNotifyPropsCache.Name())
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelByNameCache.Name())
	}
}

func (s LocalCacheChannelStore) InvalidatePinnedPostCount(channelId string) {
	s.rootStore.doInvalidateCacheCluster(s.rootStore.channelPinnedPostCountsCache, channelId, nil)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelPinnedPostCountsCache.Name())
	}
}

func (s LocalCacheChannelStore) InvalidateMemberCount(channelId string) {
	s.rootStore.doInvalidateCacheCluster(s.rootStore.channelMemberCountsCache, channelId, nil)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelMemberCountsCache.Name())
	}
}

func (s LocalCacheChannelStore) InvalidateGuestCount(channelId string) {
	s.rootStore.doInvalidateCacheCluster(s.rootStore.channelGuestCountCache, channelId, nil)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelGuestCountCache.Name())
	}
}

func (s LocalCacheChannelStore) InvalidateChannel(channelId string) {
	s.rootStore.doInvalidateCacheCluster(s.rootStore.channelByIdCache, channelId, nil)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelByIdCache.Name())
	}
}

func (s LocalCacheChannelStore) InvalidateAllChannelMembersForUser(userId string) {
	s.rootStore.doInvalidateCacheCluster(s.rootStore.channelMembersForUserCache, userId, nil)
	s.rootStore.doInvalidateCacheCluster(s.rootStore.channelMembersForUserCache, userId+"_deleted", nil)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelMembersForUserCache.Name())
	}
}

func (s LocalCacheChannelStore) InvalidateCacheForChannelMembersNotifyProps(channelId string) {
	s.rootStore.doInvalidateCacheCluster(s.rootStore.channelMembersNotifyPropsCache, channelId, nil)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelMembersNotifyPropsCache.Name())
	}
}

func (s LocalCacheChannelStore) InvalidateChannelByName(teamId, name string) {
	props := make(map[string]string)
	props["name"] = name
	if teamId == "" {
		props["id"] = "dm"
	} else {
		props["id"] = teamId
	}

	s.rootStore.doInvalidateCacheCluster(s.rootStore.channelByNameCache, teamId+name, props)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter(s.rootStore.channelByNameCache.Name())
	}
}

func (s LocalCacheChannelStore) GetMemberCount(channelId string, allowFromCache bool) (int64, error) {
	if allowFromCache {
		var count int64
		if err := s.rootStore.doStandardReadCache(s.rootStore.channelMemberCountsCache, channelId, &count); err == nil {
			return count, nil
		}
	}
	count, err := s.ChannelStore.GetMemberCount(channelId, allowFromCache)

	if allowFromCache && err == nil {
		s.rootStore.doStandardAddToCache(s.rootStore.channelMemberCountsCache, channelId, count)
	}

	return count, err
}

func (s LocalCacheChannelStore) GetGuestCount(channelId string, allowFromCache bool) (int64, error) {
	if allowFromCache {
		var count int64
		if err := s.rootStore.doStandardReadCache(s.rootStore.channelGuestCountCache, channelId, &count); err == nil {
			return count, nil
		}
	}
	count, err := s.ChannelStore.GetGuestCount(channelId, allowFromCache)

	if allowFromCache && err == nil {
		s.rootStore.doStandardAddToCache(s.rootStore.channelGuestCountCache, channelId, count)
	}

	return count, err
}

func (s LocalCacheChannelStore) GetMemberCountFromCache(channelId string) int64 {
	var count int64
	if err := s.rootStore.doStandardReadCache(s.rootStore.channelMemberCountsCache, channelId, &count); err == nil {
		return count
	}

	count, err := s.GetMemberCount(channelId, true)
	if err != nil {
		return 0
	}

	return count
}

func (s LocalCacheChannelStore) GetPinnedPostCount(channelId string, allowFromCache bool) (int64, error) {
	if allowFromCache {
		var count int64
		if err := s.rootStore.doStandardReadCache(s.rootStore.channelPinnedPostCountsCache, channelId, &count); err == nil {
			return count, nil
		}
	}

	count, err := s.ChannelStore.GetPinnedPostCount(channelId, allowFromCache)

	if err != nil {
		return 0, err
	}

	if allowFromCache {
		s.rootStore.doStandardAddToCache(s.rootStore.channelPinnedPostCountsCache, channelId, count)
	}

	return count, nil
}

func (s LocalCacheChannelStore) Get(id string, allowFromCache bool) (*model.Channel, error) {
	if allowFromCache {
		var cacheItem *model.Channel
		if err := s.rootStore.doStandardReadCache(s.rootStore.channelByIdCache, id, &cacheItem); err == nil {
			return cacheItem, nil
		}
	}

	ch, err := s.ChannelStore.Get(id, allowFromCache)

	if allowFromCache && err == nil {
		s.rootStore.doStandardAddToCache(s.rootStore.channelByIdCache, id, ch)
	}

	return ch, err
}

func (s LocalCacheChannelStore) GetMany(ids []string, allowFromCache bool) (model.ChannelList, error) {
	var foundChannels []*model.Channel
	var channelsToQuery []string

	if allowFromCache {
		for _, id := range ids {
			var ch *model.Channel
			if err := s.rootStore.doStandardReadCache(s.rootStore.channelByIdCache, id, &ch); err == nil {
				foundChannels = append(foundChannels, ch)
			} else {
				channelsToQuery = append(channelsToQuery, id)
			}
		}
	}

	if channelsToQuery == nil {
		return foundChannels, nil
	}

	channels, err := s.ChannelStore.GetMany(channelsToQuery, allowFromCache)
	if err != nil {
		return nil, err
	}

	for _, ch := range channels {
		s.rootStore.doStandardAddToCache(s.rootStore.channelByIdCache, ch.Id, ch)
	}

	return append(foundChannels, channels...), nil
}

func (s LocalCacheChannelStore) GetAllChannelMembersForUser(ctx request.CTX, userId string, allowFromCache bool, includeDeleted bool) (map[string]*model.AllChannelMember, error) {
	cache_key := userId
	if includeDeleted {
		cache_key += "_deleted"
	}
	if allowFromCache {
		ids := make(map[string]*model.AllChannelMember)
		if err := s.rootStore.doStandardReadCache(s.rootStore.channelMembersForUserCache, cache_key, &ids); err == nil {
			return ids, nil
		}
	}

	ids, err := s.ChannelStore.GetAllChannelMembersForUser(ctx, userId, allowFromCache, includeDeleted)
	if err != nil {
		return nil, err
	}

	if allowFromCache {
		s.rootStore.doStandardAddToCache(s.rootStore.channelMembersForUserCache, cache_key, ids)
	}

	return ids, nil
}

func (s LocalCacheChannelStore) GetAllChannelMembersNotifyPropsForChannel(channelId string, allowFromCache bool) (map[string]model.StringMap, error) {
	if allowFromCache {
		var cacheItem map[string]model.StringMap
		if err := s.rootStore.doStandardReadCache(s.rootStore.channelMembersNotifyPropsCache, channelId, &cacheItem); err == nil {
			return cacheItem, nil
		}
	}

	props, err := s.ChannelStore.GetAllChannelMembersNotifyPropsForChannel(channelId, allowFromCache)
	if err != nil {
		return nil, err
	}

	if allowFromCache {
		s.rootStore.doStandardAddToCache(s.rootStore.channelMembersNotifyPropsCache, channelId, props)
	}

	return props, nil
}

func (s LocalCacheChannelStore) GetByNamesIncludeDeleted(teamId string, names []string, allowFromCache bool) ([]*model.Channel, error) {
	return s.getByNames(teamId, names, allowFromCache, true)
}

func (s LocalCacheChannelStore) GetByNames(teamId string, names []string, allowFromCache bool) ([]*model.Channel, error) {
	return s.getByNames(teamId, names, allowFromCache, false)
}

func (s LocalCacheChannelStore) getByNames(teamId string, names []string, allowFromCache, includeArchivedChannels bool) ([]*model.Channel, error) {
	var channels []*model.Channel

	if allowFromCache {
		var misses []string
		visited := make(map[string]struct{})
		for _, name := range names {
			if _, ok := visited[name]; ok {
				continue
			}
			visited[name] = struct{}{}
			var cacheItem *model.Channel

			if err := s.rootStore.doStandardReadCache(s.rootStore.channelByNameCache, teamId+name, &cacheItem); err == nil {
				if includeArchivedChannels || cacheItem.DeleteAt == 0 {
					channels = append(channels, cacheItem)
				}
			} else {
				misses = append(misses, name)
			}
		}
		names = misses
	}

	if len(names) > 0 {
		var dbChannels []*model.Channel
		var err error
		if includeArchivedChannels {
			dbChannels, err = s.ChannelStore.GetByNamesIncludeDeleted(teamId, names, allowFromCache)
		} else {
			dbChannels, err = s.ChannelStore.GetByNames(teamId, names, allowFromCache)
		}
		if err != nil {
			return nil, err
		}

		for _, channel := range dbChannels {
			if allowFromCache {
				s.rootStore.doStandardAddToCache(s.rootStore.channelByNameCache, teamId+channel.Name, channel)
			}
			channels = append(channels, channel)
		}
	}

	return channels, nil
}

func (s LocalCacheChannelStore) GetByNameIncludeDeleted(teamId string, name string, allowFromCache bool) (*model.Channel, error) {
	return s.getByName(teamId, name, allowFromCache, true)
}

func (s LocalCacheChannelStore) GetByName(teamId string, name string, allowFromCache bool) (*model.Channel, error) {
	return s.getByName(teamId, name, allowFromCache, false)
}

func (s LocalCacheChannelStore) getByName(teamId string, name string, allowFromCache, includeArchivedChannels bool) (*model.Channel, error) {
	var channel *model.Channel

	if allowFromCache {
		if err := s.rootStore.doStandardReadCache(s.rootStore.channelByNameCache, teamId+name, &channel); err == nil {
			if includeArchivedChannels || channel.DeleteAt == 0 {
				return channel, nil
			}
		}
	}

	var err error
	if includeArchivedChannels {
		channel, err = s.ChannelStore.GetByNameIncludeDeleted(teamId, name, allowFromCache)
	} else {
		channel, err = s.ChannelStore.GetByName(teamId, name, allowFromCache)
	}
	if err != nil {
		return nil, err
	}

	if allowFromCache {
		s.rootStore.doStandardAddToCache(s.rootStore.channelByNameCache, teamId+name, channel)
	}

	return channel, nil
}

func (s LocalCacheChannelStore) SaveMember(rctx request.CTX, member *model.ChannelMember) (*model.ChannelMember, error) {
	member, err := s.ChannelStore.SaveMember(rctx, member)
	if err != nil {
		return nil, err
	}
	s.InvalidateMemberCount(member.ChannelId)
	return member, nil
}

func (s LocalCacheChannelStore) SaveMultipleMembers(members []*model.ChannelMember) ([]*model.ChannelMember, error) {
	members, err := s.ChannelStore.SaveMultipleMembers(members)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		s.InvalidateMemberCount(member.ChannelId)
	}
	return members, nil
}

func (s LocalCacheChannelStore) GetChannelsMemberCount(channelIDs []string) (_ map[string]int64, err error) {
	counts := make(map[string]int64)
	remainingChannels := make([]string, 0)

	for _, channelID := range channelIDs {
		var cacheItem int64
		err := s.rootStore.doStandardReadCache(s.rootStore.channelMemberCountsCache, channelID, &cacheItem)
		if err == nil {
			counts[channelID] = cacheItem
		} else {
			remainingChannels = append(remainingChannels, channelID)
		}
	}

	if len(remainingChannels) > 0 {
		remainingChannels, err := s.ChannelStore.GetChannelsMemberCount(remainingChannels)
		if err != nil {
			return nil, err
		}

		for id, count := range remainingChannels {
			s.rootStore.doStandardAddToCache(s.rootStore.channelMemberCountsCache, id, count)
			counts[id] = count
		}
	}

	return counts, nil
}

func (s LocalCacheChannelStore) UpdateMember(rctx request.CTX, member *model.ChannelMember) (*model.ChannelMember, error) {
	member, err := s.ChannelStore.UpdateMember(rctx, member)
	if err != nil {
		return nil, err
	}
	s.InvalidateMemberCount(member.ChannelId)
	return member, nil
}

func (s LocalCacheChannelStore) UpdateMultipleMembers(members []*model.ChannelMember) ([]*model.ChannelMember, error) {
	members, err := s.ChannelStore.UpdateMultipleMembers(members)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		s.InvalidateMemberCount(member.ChannelId)
	}
	return members, nil
}

func (s LocalCacheChannelStore) RemoveMember(rctx request.CTX, channelId, userId string) error {
	err := s.ChannelStore.RemoveMember(rctx, channelId, userId)
	if err != nil {
		return err
	}
	s.InvalidateMemberCount(channelId)
	return nil
}

func (s LocalCacheChannelStore) RemoveMembers(rctx request.CTX, channelId string, userIds []string) error {
	err := s.ChannelStore.RemoveMembers(rctx, channelId, userIds)
	if err != nil {
		return err
	}
	s.InvalidateMemberCount(channelId)
	return nil
}

func (s LocalCacheChannelStore) GetCachedAllChannelMembersForUser(userId string, includeDeleted bool) (map[string]*model.AllChannelMember, error) {
	cache_key := userId
	if includeDeleted {
		cache_key += "_deleted"
	}

	ids := make(map[string]*model.AllChannelMember)
	if err := s.rootStore.doStandardReadCache(s.rootStore.channelMembersForUserCache, cache_key, &ids); err != nil {
		if !errors.Is(err, cache.ErrKeyNotFound) {
			return nil, err
		}
	}

	return ids, nil
}
