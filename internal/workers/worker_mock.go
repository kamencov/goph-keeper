// Code generated by MockGen. DO NOT EDIT.
// Source: worker.go

// Package workers is a generated GoMock package.
package workers

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
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

// PushData mocks base method.
func (m *Mockservice) PushData(ctx context.Context, conn *grpc.ClientConn) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushData", ctx, conn)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushData indicates an expected call of PushData.
func (mr *MockserviceMockRecorder) PushData(ctx, conn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushData", reflect.TypeOf((*Mockservice)(nil).PushData), ctx, conn)
}
