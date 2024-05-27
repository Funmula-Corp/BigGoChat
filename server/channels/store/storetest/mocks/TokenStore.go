// Code generated by mockery v2.42.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	mock "github.com/stretchr/testify/mock"
)

// TokenStore is an autogenerated mock type for the TokenStore type
type TokenStore struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields: expiryTime
func (_m *TokenStore) Cleanup(expiryTime int64) {
	_m.Called(expiryTime)
}

// Delete provides a mock function with given fields: token
func (_m *TokenStore) Delete(token string) error {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllTokensByType provides a mock function with given fields: tokenType
func (_m *TokenStore) GetAllTokensByType(tokenType string) ([]*model.Token, error) {
	ret := _m.Called(tokenType)

	if len(ret) == 0 {
		panic("no return value specified for GetAllTokensByType")
	}

	var r0 []*model.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]*model.Token, error)); ok {
		return rf(tokenType)
	}
	if rf, ok := ret.Get(0).(func(string) []*model.Token); ok {
		r0 = rf(tokenType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByToken provides a mock function with given fields: token
func (_m *TokenStore) GetByToken(token string) (*model.Token, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for GetByToken")
	}

	var r0 *model.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Token, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Token); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveAllTokensByType provides a mock function with given fields: tokenType
func (_m *TokenStore) RemoveAllTokensByType(tokenType string) error {
	ret := _m.Called(tokenType)

	if len(ret) == 0 {
		panic("no return value specified for RemoveAllTokensByType")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(tokenType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: recovery
func (_m *TokenStore) Save(recovery *model.Token) error {
	ret := _m.Called(recovery)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Token) error); ok {
		r0 = rf(recovery)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTokenStore creates a new instance of TokenStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenStore {
	mock := &TokenStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
