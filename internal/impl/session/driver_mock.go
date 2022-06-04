// Code generated by MockGen. DO NOT EDIT.
// Source: session_driver.go

// Package session is a generated GoMock package.
package session

import (
	reflect "reflect"

	domain "github.com/boris-army/server/internal/core/domain"
	ports "github.com/boris-army/server/internal/core/ports"
	gomock "github.com/golang/mock/gomock"
)

// MockSessionDriver is a mock of SessionDriver interface.
type MockSessionDriver struct {
	ctrl     *gomock.Controller
	recorder *MockSessionDriverMockRecorder
}

// MockSessionDriverMockRecorder is the mock recorder for MockSessionDriver.
type MockSessionDriverMockRecorder struct {
	mock *MockSessionDriver
}

// NewMockSessionDriver creates a new mock instance.
func NewMockSessionDriver(ctrl *gomock.Controller) *MockSessionDriver {
	mock := &MockSessionDriver{ctrl: ctrl}
	mock.recorder = &MockSessionDriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionDriver) EXPECT() *MockSessionDriverMockRecorder {
	return m.recorder
}

// CreateHttp mocks base method.
func (m *MockSessionDriver) CreateHttp(create *ports.CommandSessionHttpCreate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHttp", create)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateHttp indicates an expected call of CreateHttp.
func (mr *MockSessionDriverMockRecorder) CreateHttp(create interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHttp", reflect.TypeOf((*MockSessionDriver)(nil).CreateHttp), create)
}

// DecodeHttpTokenTo mocks base method.
func (m *MockSessionDriver) DecodeHttpTokenTo(dst *domain.SessionHttpToken, src []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodeHttpTokenTo", dst, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecodeHttpTokenTo indicates an expected call of DecodeHttpTokenTo.
func (mr *MockSessionDriverMockRecorder) DecodeHttpTokenTo(dst, src interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodeHttpTokenTo", reflect.TypeOf((*MockSessionDriver)(nil).DecodeHttpTokenTo), dst, src)
}
