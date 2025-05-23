// Code generated by MockGen. DO NOT EDIT.
// Source: git.biggo.com/Funmula/BigGoChat/server/public/pluginapi/experimental/panel (interfaces: Store)

// Package mock_panel is a generated GoMock package.
package mock_panel

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// DeletePanelPostID mocks base method.
func (m *MockStore) DeletePanelPostID(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePanelPostID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePanelPostID indicates an expected call of DeletePanelPostID.
func (mr *MockStoreMockRecorder) DeletePanelPostID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePanelPostID", reflect.TypeOf((*MockStore)(nil).DeletePanelPostID), arg0)
}

// GetPanelPostID mocks base method.
func (m *MockStore) GetPanelPostID(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPanelPostID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPanelPostID indicates an expected call of GetPanelPostID.
func (mr *MockStoreMockRecorder) GetPanelPostID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPanelPostID", reflect.TypeOf((*MockStore)(nil).GetPanelPostID), arg0)
}

// SetPanelPostID mocks base method.
func (m *MockStore) SetPanelPostID(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPanelPostID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPanelPostID indicates an expected call of SetPanelPostID.
func (mr *MockStoreMockRecorder) SetPanelPostID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPanelPostID", reflect.TypeOf((*MockStore)(nil).SetPanelPostID), arg0, arg1)
}
