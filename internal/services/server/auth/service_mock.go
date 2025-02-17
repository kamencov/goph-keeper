// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package auth is a generated GoMock package.
package auth

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockstorageAuth is a mock of storageAuth interface.
type MockstorageAuth struct {
	ctrl     *gomock.Controller
	recorder *MockstorageAuthMockRecorder
}

// MockstorageAuthMockRecorder is the mock recorder for MockstorageAuth.
type MockstorageAuthMockRecorder struct {
	mock *MockstorageAuth
}

// NewMockstorageAuth creates a new mock instance.
func NewMockstorageAuth(ctrl *gomock.Controller) *MockstorageAuth {
	mock := &MockstorageAuth{ctrl: ctrl}
	mock.recorder = &MockstorageAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockstorageAuth) EXPECT() *MockstorageAuthMockRecorder {
	return m.recorder
}

// CheckPassword mocks base method.
func (m *MockstorageAuth) CheckPassword(login string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPassword", login)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// CheckPassword indicates an expected call of CheckPassword.
func (mr *MockstorageAuthMockRecorder) CheckPassword(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPassword", reflect.TypeOf((*MockstorageAuth)(nil).CheckPassword), login)
}

// CheckUser mocks base method.
func (m *MockstorageAuth) CheckUser(ctx context.Context, login string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", ctx, login)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockstorageAuthMockRecorder) CheckUser(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockstorageAuth)(nil).CheckUser), ctx, login)
}

// GetUserIDByLogin mocks base method.
func (m *MockstorageAuth) GetUserIDByLogin(ctx context.Context, login string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDByLogin", ctx, login)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDByLogin indicates an expected call of GetUserIDByLogin.
func (mr *MockstorageAuthMockRecorder) GetUserIDByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDByLogin", reflect.TypeOf((*MockstorageAuth)(nil).GetUserIDByLogin), ctx, login)
}

// GetUserIDByToken mocks base method.
func (m *MockstorageAuth) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDByToken", ctx, token)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDByToken indicates an expected call of GetUserIDByToken.
func (mr *MockstorageAuthMockRecorder) GetUserIDByToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDByToken", reflect.TypeOf((*MockstorageAuth)(nil).GetUserIDByToken), ctx, token)
}

// SaveTableUserAndUpdateToken mocks base method.
func (m *MockstorageAuth) SaveTableUserAndUpdateToken(login, accessToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTableUserAndUpdateToken", login, accessToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveTableUserAndUpdateToken indicates an expected call of SaveTableUserAndUpdateToken.
func (mr *MockstorageAuthMockRecorder) SaveTableUserAndUpdateToken(login, accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTableUserAndUpdateToken", reflect.TypeOf((*MockstorageAuth)(nil).SaveTableUserAndUpdateToken), login, accessToken)
}

// SaveUser mocks base method.
func (m *MockstorageAuth) SaveUser(ctx context.Context, login, hashPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", ctx, login, hashPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUser indicates an expected call of SaveUser.
func (mr *MockstorageAuthMockRecorder) SaveUser(ctx, login, hashPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUser", reflect.TypeOf((*MockstorageAuth)(nil).SaveUser), ctx, login, hashPassword)
}
