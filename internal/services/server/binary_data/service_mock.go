// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package binary_data is a generated GoMock package.
package binary_data

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Mockstorage is a mock of storage interface.
type Mockstorage struct {
	ctrl     *gomock.Controller
	recorder *MockstorageMockRecorder
}

// MockstorageMockRecorder is the mock recorder for Mockstorage.
type MockstorageMockRecorder struct {
	mock *Mockstorage
}

// NewMockstorage creates a new mock instance.
func NewMockstorage(ctrl *gomock.Controller) *Mockstorage {
	mock := &Mockstorage{ctrl: ctrl}
	mock.recorder = &MockstorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockstorage) EXPECT() *MockstorageMockRecorder {
	return m.recorder
}

// DeletedBinary mocks base method.
func (m *Mockstorage) DeletedBinary(ctx context.Context, userID int, data string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletedBinary", ctx, userID, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletedBinary indicates an expected call of DeletedBinary.
func (mr *MockstorageMockRecorder) DeletedBinary(ctx, userID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletedBinary", reflect.TypeOf((*Mockstorage)(nil).DeletedBinary), ctx, userID, data)
}

// GetUserIDByToken mocks base method.
func (m *Mockstorage) GetUserIDByToken(ctx context.Context, accessToken string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDByToken", ctx, accessToken)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDByToken indicates an expected call of GetUserIDByToken.
func (mr *MockstorageMockRecorder) GetUserIDByToken(ctx, accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDByToken", reflect.TypeOf((*Mockstorage)(nil).GetUserIDByToken), ctx, accessToken)
}

// SaveBinaryDataBinary mocks base method.
func (m *Mockstorage) SaveBinaryDataBinary(ctx context.Context, uid int, data string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBinaryDataBinary", ctx, uid, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBinaryDataBinary indicates an expected call of SaveBinaryDataBinary.
func (mr *MockstorageMockRecorder) SaveBinaryDataBinary(ctx, uid, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBinaryDataBinary", reflect.TypeOf((*Mockstorage)(nil).SaveBinaryDataBinary), ctx, uid, data)
}
