// Code generated by MockGen. DO NOT EDIT.
// Source: user_repository.go

// Package user is a generated GoMock package.
package user

import (
	reflect "reflect"

	domain "github.com/boris-army/server/internal/core/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryUser is a mock of RepositoryUser interface.
type MockRepositoryUser struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryUserMockRecorder
}

// MockRepositoryUserMockRecorder is the mock recorder for MockRepositoryUser.
type MockRepositoryUserMockRecorder struct {
	mock *MockRepositoryUser
}

// NewMockRepositoryUser creates a new mock instance.
func NewMockRepositoryUser(ctrl *gomock.Controller) *MockRepositoryUser {
	mock := &MockRepositoryUser{ctrl: ctrl}
	mock.recorder = &MockRepositoryUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryUser) EXPECT() *MockRepositoryUserMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepositoryUser) Create(arg0 *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryUserMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepositoryUser)(nil).Create), arg0)
}
