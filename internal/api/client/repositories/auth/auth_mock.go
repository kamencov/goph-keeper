// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package auth is a generated GoMock package.
package auth

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Mockservice is a mock of service interface.
type Mockservice struct {
	ctrl     *gomock.Controller
	recorder *MockserviceMockRecorder
}

// MockserviceMockRecorder is the mock recorder for Mockservice.
type MockserviceMockRecorder struct {
	mock *Mockservice
}

// NewMockservice creates a new mock instance.
func NewMockservice(ctrl *gomock.Controller) *Mockservice {
	mock := &Mockservice{ctrl: ctrl}
	mock.recorder = &MockserviceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockservice) EXPECT() *MockserviceMockRecorder {
	return m.recorder
}

// CheckUser mocks base method.
func (m *Mockservice) CheckUser(ctx context.Context, login, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", ctx, login, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockserviceMockRecorder) CheckUser(ctx, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*Mockservice)(nil).CheckUser), ctx, login, password)
}

// SaveTokenInBase mocks base method.
func (m *Mockservice) SaveTokenInBase(ctx context.Context, login, password, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTokenInBase", ctx, login, password, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveTokenInBase indicates an expected call of SaveTokenInBase.
func (mr *MockserviceMockRecorder) SaveTokenInBase(ctx, login, password, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTokenInBase", reflect.TypeOf((*Mockservice)(nil).SaveTokenInBase), ctx, login, password, token)
}
