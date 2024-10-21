// Code generated by MockGen. DO NOT EDIT.
// Source: git.biggo.com/Funmula/BigGoChat/server/public/pluginapi/experimental/bot (interfaces: Bot)

// Package mock_bot is a generated GoMock package.
package mock_bot

import (
	reflect "reflect"

	model "git.biggo.com/Funmula/BigGoChat/server/public/model"
	gomock "github.com/golang/mock/gomock"
)

// MockBot is a mock of Bot interface.
type MockBot struct {
	ctrl     *gomock.Controller
	recorder *MockBotMockRecorder
}

// MockBotMockRecorder is the mock recorder for MockBot.
type MockBotMockRecorder struct {
	mock *MockBot
}

// NewMockBot creates a new mock instance.
func NewMockBot(ctrl *gomock.Controller) *MockBot {
	mock := &MockBot{ctrl: ctrl}
	mock.recorder = &MockBotMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBot) EXPECT() *MockBotMockRecorder {
	return m.recorder
}

// Ensure mocks base method.
func (m *MockBot) Ensure(arg0 *model.Bot, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ensure", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ensure indicates an expected call of Ensure.
func (mr *MockBotMockRecorder) Ensure(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ensure", reflect.TypeOf((*MockBot)(nil).Ensure), arg0, arg1)
}

// MattermostUserID mocks base method.
func (m *MockBot) MattermostUserID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MattermostUserID")
	ret0, _ := ret[0].(string)
	return ret0
}

// MattermostUserID indicates an expected call of MattermostUserID.
func (mr *MockBotMockRecorder) MattermostUserID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MattermostUserID", reflect.TypeOf((*MockBot)(nil).MattermostUserID))
}

// String mocks base method.
func (m *MockBot) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockBotMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockBot)(nil).String))
}
