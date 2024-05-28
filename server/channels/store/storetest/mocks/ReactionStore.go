// Code generated by mockery v2.42.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	mock "github.com/stretchr/testify/mock"
)

// ReactionStore is an autogenerated mock type for the ReactionStore type
type ReactionStore struct {
	mock.Mock
}

// BulkGetForPosts provides a mock function with given fields: postIds
func (_m *ReactionStore) BulkGetForPosts(postIds []string) ([]*model.Reaction, error) {
	ret := _m.Called(postIds)

	if len(ret) == 0 {
		panic("no return value specified for BulkGetForPosts")
	}

	var r0 []*model.Reaction
	var r1 error
	if rf, ok := ret.Get(0).(func([]string) ([]*model.Reaction, error)); ok {
		return rf(postIds)
	}
	if rf, ok := ret.Get(0).(func([]string) []*model.Reaction); ok {
		r0 = rf(postIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Reaction)
		}
	}

	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(postIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: reaction
func (_m *ReactionStore) Delete(reaction *model.Reaction) (*model.Reaction, error) {
	ret := _m.Called(reaction)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 *model.Reaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Reaction) (*model.Reaction, error)); ok {
		return rf(reaction)
	}
	if rf, ok := ret.Get(0).(func(*model.Reaction) *model.Reaction); ok {
		r0 = rf(reaction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Reaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Reaction) error); ok {
		r1 = rf(reaction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAllWithEmojiName provides a mock function with given fields: emojiName
func (_m *ReactionStore) DeleteAllWithEmojiName(emojiName string) error {
	ret := _m.Called(emojiName)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAllWithEmojiName")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(emojiName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOrphanedRowsByIds provides a mock function with given fields: r
func (_m *ReactionStore) DeleteOrphanedRowsByIds(r *model.RetentionIdsForDeletion) error {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOrphanedRowsByIds")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.RetentionIdsForDeletion) error); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExistsOnPost provides a mock function with given fields: postId, emojiName
func (_m *ReactionStore) ExistsOnPost(postId string, emojiName string) (bool, error) {
	ret := _m.Called(postId, emojiName)

	if len(ret) == 0 {
		panic("no return value specified for ExistsOnPost")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (bool, error)); ok {
		return rf(postId, emojiName)
	}
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(postId, emojiName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(postId, emojiName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForPost provides a mock function with given fields: postID, allowFromCache
func (_m *ReactionStore) GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error) {
	ret := _m.Called(postID, allowFromCache)

	if len(ret) == 0 {
		panic("no return value specified for GetForPost")
	}

	var r0 []*model.Reaction
	var r1 error
	if rf, ok := ret.Get(0).(func(string, bool) ([]*model.Reaction, error)); ok {
		return rf(postID, allowFromCache)
	}
	if rf, ok := ret.Get(0).(func(string, bool) []*model.Reaction); ok {
		r0 = rf(postID, allowFromCache)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Reaction)
		}
	}

	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(postID, allowFromCache)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForPostSince provides a mock function with given fields: postId, since, excludeRemoteId, inclDeleted
func (_m *ReactionStore) GetForPostSince(postId string, since int64, excludeRemoteId string, inclDeleted bool) ([]*model.Reaction, error) {
	ret := _m.Called(postId, since, excludeRemoteId, inclDeleted)

	if len(ret) == 0 {
		panic("no return value specified for GetForPostSince")
	}

	var r0 []*model.Reaction
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int64, string, bool) ([]*model.Reaction, error)); ok {
		return rf(postId, since, excludeRemoteId, inclDeleted)
	}
	if rf, ok := ret.Get(0).(func(string, int64, string, bool) []*model.Reaction); ok {
		r0 = rf(postId, since, excludeRemoteId, inclDeleted)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Reaction)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int64, string, bool) error); ok {
		r1 = rf(postId, since, excludeRemoteId, inclDeleted)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUniqueCountForPost provides a mock function with given fields: postId
func (_m *ReactionStore) GetUniqueCountForPost(postId string) (int, error) {
	ret := _m.Called(postId)

	if len(ret) == 0 {
		panic("no return value specified for GetUniqueCountForPost")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int, error)); ok {
		return rf(postId)
	}
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(postId)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(postId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PermanentDeleteBatch provides a mock function with given fields: endTime, limit
func (_m *ReactionStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	ret := _m.Called(endTime, limit)

	if len(ret) == 0 {
		panic("no return value specified for PermanentDeleteBatch")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int64) (int64, error)); ok {
		return rf(endTime, limit)
	}
	if rf, ok := ret.Get(0).(func(int64, int64) int64); ok {
		r0 = rf(endTime, limit)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(endTime, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PermanentDeleteByUser provides a mock function with given fields: userID
func (_m *ReactionStore) PermanentDeleteByUser(userID string) error {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for PermanentDeleteByUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: reaction
func (_m *ReactionStore) Save(reaction *model.Reaction) (*model.Reaction, error) {
	ret := _m.Called(reaction)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 *model.Reaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Reaction) (*model.Reaction, error)); ok {
		return rf(reaction)
	}
	if rf, ok := ret.Get(0).(func(*model.Reaction) *model.Reaction); ok {
		r0 = rf(reaction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Reaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Reaction) error); ok {
		r1 = rf(reaction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewReactionStore creates a new instance of ReactionStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReactionStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReactionStore {
	mock := &ReactionStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
