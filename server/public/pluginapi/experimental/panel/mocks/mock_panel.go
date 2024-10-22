// Code generated by MockGen. DO NOT EDIT.
// Source: git.biggo.com/Funmula/BigGoChat/server/public/pluginapi/experimental/panel (interfaces: Panel)

// Package mock_panel is a generated GoMock package.
package mock_panel

import (
	reflect "reflect"

	model "git.biggo.com/Funmula/BigGoChat/server/public/model"
	gomock "github.com/golang/mock/gomock"
)

// MockPanel is a mock of Panel interface.
type MockPanel struct {
	ctrl     *gomock.Controller
	recorder *MockPanelMockRecorder
}

// MockPanelMockRecorder is the mock recorder for MockPanel.
type MockPanelMockRecorder struct {
	mock *MockPanel
}

// NewMockPanel creates a new mock instance.
func NewMockPanel(ctrl *gomock.Controller) *MockPanel {
	mock := &MockPanel{ctrl: ctrl}
	mock.recorder = &MockPanelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPanel) EXPECT() *MockPanelMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockPanel) Clear(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Clear indicates an expected call of Clear.
func (mr *MockPanelMockRecorder) Clear(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockPanel)(nil).Clear), arg0)
}

// GetSettingIDs mocks base method.
func (m *MockPanel) GetSettingIDs() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSettingIDs")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetSettingIDs indicates an expected call of GetSettingIDs.
func (mr *MockPanelMockRecorder) GetSettingIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSettingIDs", reflect.TypeOf((*MockPanel)(nil).GetSettingIDs))
}

// Print mocks base method.
func (m *MockPanel) Print(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Print", arg0)
}

// Print indicates an expected call of Print.
func (mr *MockPanelMockRecorder) Print(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Print", reflect.TypeOf((*MockPanel)(nil).Print), arg0)
}

// Set mocks base method.
func (m *MockPanel) Set(arg0, arg1 string, arg2 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockPanelMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockPanel)(nil).Set), arg0, arg1, arg2)
}

// ToPost mocks base method.
func (m *MockPanel) ToPost(arg0 string) (*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToPost", arg0)
	ret0, _ := ret[0].(*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToPost indicates an expected call of ToPost.
func (mr *MockPanelMockRecorder) ToPost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToPost", reflect.TypeOf((*MockPanel)(nil).ToPost), arg0)
}

// URL mocks base method.
func (m *MockPanel) URL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "URL")
	ret0, _ := ret[0].(string)
	return ret0
}

// URL indicates an expected call of URL.
func (mr *MockPanelMockRecorder) URL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URL", reflect.TypeOf((*MockPanel)(nil).URL))
}
