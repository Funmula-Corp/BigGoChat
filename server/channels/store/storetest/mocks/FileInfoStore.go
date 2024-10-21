// Code generated by mockery v2.42.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "git.biggo.com/Funmula/BigGoChat/server/public/model"
	request "git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
	mock "github.com/stretchr/testify/mock"
)

// FileInfoStore is an autogenerated mock type for the FileInfoStore type
type FileInfoStore struct {
	mock.Mock
}

// AttachToPost provides a mock function with given fields: c, fileID, postID, channelID, creatorID
func (_m *FileInfoStore) AttachToPost(c request.CTX, fileID string, postID string, channelID string, creatorID string) error {
	ret := _m.Called(c, fileID, postID, channelID, creatorID)

	if len(ret) == 0 {
		panic("no return value specified for AttachToPost")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(request.CTX, string, string, string, string) error); ok {
		r0 = rf(c, fileID, postID, channelID, creatorID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearCaches provides a mock function with given fields:
func (_m *FileInfoStore) ClearCaches() {
	_m.Called()
}

// CountAll provides a mock function with given fields:
func (_m *FileInfoStore) CountAll() (int64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CountAll")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteForPost provides a mock function with given fields: c, postID
func (_m *FileInfoStore) DeleteForPost(c request.CTX, postID string) (string, error) {
	ret := _m.Called(c, postID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteForPost")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, string) (string, error)); ok {
		return rf(c, postID)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, string) string); ok {
		r0 = rf(c, postID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(request.CTX, string) error); ok {
		r1 = rf(c, postID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *FileInfoStore) Get(id string) (*model.FileInfo, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.FileInfo, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.FileInfo); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByIds provides a mock function with given fields: ids
func (_m *FileInfoStore) GetByIds(ids []string) ([]*model.FileInfo, error) {
	ret := _m.Called(ids)

	if len(ret) == 0 {
		panic("no return value specified for GetByIds")
	}

	var r0 []*model.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func([]string) ([]*model.FileInfo, error)); ok {
		return rf(ids)
	}
	if rf, ok := ret.Get(0).(func([]string) []*model.FileInfo); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByPath provides a mock function with given fields: path
func (_m *FileInfoStore) GetByPath(path string) (*model.FileInfo, error) {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for GetByPath")
	}

	var r0 *model.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.FileInfo, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) *model.FileInfo); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFilesBatchForIndexing provides a mock function with given fields: startTime, startFileID, includeDeleted, limit
func (_m *FileInfoStore) GetFilesBatchForIndexing(startTime int64, startFileID string, includeDeleted bool, limit int) ([]*model.FileForIndexing, error) {
	ret := _m.Called(startTime, startFileID, includeDeleted, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetFilesBatchForIndexing")
	}

	var r0 []*model.FileForIndexing
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, string, bool, int) ([]*model.FileForIndexing, error)); ok {
		return rf(startTime, startFileID, includeDeleted, limit)
	}
	if rf, ok := ret.Get(0).(func(int64, string, bool, int) []*model.FileForIndexing); ok {
		r0 = rf(startTime, startFileID, includeDeleted, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FileForIndexing)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, string, bool, int) error); ok {
		r1 = rf(startTime, startFileID, includeDeleted, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForPost provides a mock function with given fields: postID, readFromMaster, includeDeleted, allowFromCache
func (_m *FileInfoStore) GetForPost(postID string, readFromMaster bool, includeDeleted bool, allowFromCache bool) ([]*model.FileInfo, error) {
	ret := _m.Called(postID, readFromMaster, includeDeleted, allowFromCache)

	if len(ret) == 0 {
		panic("no return value specified for GetForPost")
	}

	var r0 []*model.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string, bool, bool, bool) ([]*model.FileInfo, error)); ok {
		return rf(postID, readFromMaster, includeDeleted, allowFromCache)
	}
	if rf, ok := ret.Get(0).(func(string, bool, bool, bool) []*model.FileInfo); ok {
		r0 = rf(postID, readFromMaster, includeDeleted, allowFromCache)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string, bool, bool, bool) error); ok {
		r1 = rf(postID, readFromMaster, includeDeleted, allowFromCache)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForUser provides a mock function with given fields: userID
func (_m *FileInfoStore) GetForUser(userID string) ([]*model.FileInfo, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetForUser")
	}

	var r0 []*model.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]*model.FileInfo, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) []*model.FileInfo); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFromMaster provides a mock function with given fields: id
func (_m *FileInfoStore) GetFromMaster(id string) (*model.FileInfo, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetFromMaster")
	}

	var r0 *model.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.FileInfo, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.FileInfo); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStorageUsage provides a mock function with given fields: allowFromCache, includeDeleted
func (_m *FileInfoStore) GetStorageUsage(allowFromCache bool, includeDeleted bool) (int64, error) {
	ret := _m.Called(allowFromCache, includeDeleted)

	if len(ret) == 0 {
		panic("no return value specified for GetStorageUsage")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(bool, bool) (int64, error)); ok {
		return rf(allowFromCache, includeDeleted)
	}
	if rf, ok := ret.Get(0).(func(bool, bool) int64); ok {
		r0 = rf(allowFromCache, includeDeleted)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(bool, bool) error); ok {
		r1 = rf(allowFromCache, includeDeleted)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUptoNSizeFileTime provides a mock function with given fields: n
func (_m *FileInfoStore) GetUptoNSizeFileTime(n int64) (int64, error) {
	ret := _m.Called(n)

	if len(ret) == 0 {
		panic("no return value specified for GetUptoNSizeFileTime")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (int64, error)); ok {
		return rf(n)
	}
	if rf, ok := ret.Get(0).(func(int64) int64); ok {
		r0 = rf(n)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(n)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithOptions provides a mock function with given fields: page, perPage, opt
func (_m *FileInfoStore) GetWithOptions(page int, perPage int, opt *model.GetFileInfosOptions) ([]*model.FileInfo, error) {
	ret := _m.Called(page, perPage, opt)

	if len(ret) == 0 {
		panic("no return value specified for GetWithOptions")
	}

	var r0 []*model.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, *model.GetFileInfosOptions) ([]*model.FileInfo, error)); ok {
		return rf(page, perPage, opt)
	}
	if rf, ok := ret.Get(0).(func(int, int, *model.GetFileInfosOptions) []*model.FileInfo); ok {
		r0 = rf(page, perPage, opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, *model.GetFileInfosOptions) error); ok {
		r1 = rf(page, perPage, opt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvalidateFileInfosForPostCache provides a mock function with given fields: postID, deleted
func (_m *FileInfoStore) InvalidateFileInfosForPostCache(postID string, deleted bool) {
	_m.Called(postID, deleted)
}

// PermanentDelete provides a mock function with given fields: c, fileID
func (_m *FileInfoStore) PermanentDelete(c request.CTX, fileID string) error {
	ret := _m.Called(c, fileID)

	if len(ret) == 0 {
		panic("no return value specified for PermanentDelete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(request.CTX, string) error); ok {
		r0 = rf(c, fileID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PermanentDeleteBatch provides a mock function with given fields: ctx, endTime, limit
func (_m *FileInfoStore) PermanentDeleteBatch(ctx request.CTX, endTime int64, limit int64) (int64, error) {
	ret := _m.Called(ctx, endTime, limit)

	if len(ret) == 0 {
		panic("no return value specified for PermanentDeleteBatch")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, int64, int64) (int64, error)); ok {
		return rf(ctx, endTime, limit)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, int64, int64) int64); ok {
		r0 = rf(ctx, endTime, limit)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(request.CTX, int64, int64) error); ok {
		r1 = rf(ctx, endTime, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PermanentDeleteByUser provides a mock function with given fields: ctx, userID
func (_m *FileInfoStore) PermanentDeleteByUser(ctx request.CTX, userID string) (int64, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for PermanentDeleteByUser")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, string) (int64, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, string) int64); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(request.CTX, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, info
func (_m *FileInfoStore) Save(ctx request.CTX, info *model.FileInfo) (*model.FileInfo, error) {
	ret := _m.Called(ctx, info)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 *model.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, *model.FileInfo) (*model.FileInfo, error)); ok {
		return rf(ctx, info)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, *model.FileInfo) *model.FileInfo); ok {
		r0 = rf(ctx, info)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(request.CTX, *model.FileInfo) error); ok {
		r1 = rf(ctx, info)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Search provides a mock function with given fields: ctx, paramsList, userID, teamID, page, perPage
func (_m *FileInfoStore) Search(ctx request.CTX, paramsList []*model.SearchParams, userID string, teamID string, page int, perPage int) (*model.FileInfoList, error) {
	ret := _m.Called(ctx, paramsList, userID, teamID, page, perPage)

	if len(ret) == 0 {
		panic("no return value specified for Search")
	}

	var r0 *model.FileInfoList
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, []*model.SearchParams, string, string, int, int) (*model.FileInfoList, error)); ok {
		return rf(ctx, paramsList, userID, teamID, page, perPage)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, []*model.SearchParams, string, string, int, int) *model.FileInfoList); ok {
		r0 = rf(ctx, paramsList, userID, teamID, page, perPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FileInfoList)
		}
	}

	if rf, ok := ret.Get(1).(func(request.CTX, []*model.SearchParams, string, string, int, int) error); ok {
		r1 = rf(ctx, paramsList, userID, teamID, page, perPage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetContent provides a mock function with given fields: ctx, fileID, content
func (_m *FileInfoStore) SetContent(ctx request.CTX, fileID string, content string) error {
	ret := _m.Called(ctx, fileID, content)

	if len(ret) == 0 {
		panic("no return value specified for SetContent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(request.CTX, string, string) error); ok {
		r0 = rf(ctx, fileID, content)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Upsert provides a mock function with given fields: rctx, info
func (_m *FileInfoStore) Upsert(rctx request.CTX, info *model.FileInfo) (*model.FileInfo, error) {
	ret := _m.Called(rctx, info)

	if len(ret) == 0 {
		panic("no return value specified for Upsert")
	}

	var r0 *model.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, *model.FileInfo) (*model.FileInfo, error)); ok {
		return rf(rctx, info)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, *model.FileInfo) *model.FileInfo); ok {
		r0 = rf(rctx, info)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(request.CTX, *model.FileInfo) error); ok {
		r1 = rf(rctx, info)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFileInfoStore creates a new instance of FileInfoStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFileInfoStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *FileInfoStore {
	mock := &FileInfoStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
