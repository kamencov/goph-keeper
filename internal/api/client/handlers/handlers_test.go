package handlers

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

var (
	errSave = errors.New("failed to save")
)

func TestNewHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	testServiceCredentials := NewMockserviceCredentials(ctrl)
	testServiceCards := NewMockserviceCards(ctrl)
	testServiceBinaryData := NewMockserviceBinaryData(ctrl)
	testServiceTextData := NewMockserviceTextData(ctrl)

	handler := NewHandlers(log, testServiceCredentials, testServiceTextData, testServiceBinaryData, testServiceCards)

	if handler == nil {
		t.Errorf("Handler is nil")
	}
}

func TestHandler_PostLoginAndPassword(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		resource    string
		login       string
		password    string
		expectedErr error
	}{
		{
			name:        "successful_save",
			token:       "test_token",
			resource:    "test_resource",
			login:       "test_login",
			password:    "test_password",
			expectedErr: nil,
		},
		{
			name:        "failed_null_resource",
			token:       "token",
			resource:    "",
			login:       "test_login",
			password:    "test_password",
			expectedErr: ErrNotEmpty,
		},
		{
			name:        "failed_save",
			token:       "test_token",
			resource:    "test_resource",
			login:       "test_login",
			password:    "test_password",
			expectedErr: errSave,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			testServiceCredentials := NewMockserviceCredentials(ctrl)
			testServiceCards := NewMockserviceCards(ctrl)
			testServiceBinaryData := NewMockserviceBinaryData(ctrl)
			testServiceTextData := NewMockserviceTextData(ctrl)
			testServiceCredentials.EXPECT().SaveLoginAndPassword(ctx, tt.token, tt.resource, tt.login, tt.password).Return(tt.expectedErr).AnyTimes()
			handler := NewHandlers(log, testServiceCredentials, testServiceTextData, testServiceBinaryData, testServiceCards)

			err := handler.PostLoginAndPassword(ctx, tt.token, tt.resource, tt.login, tt.password)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("PostLoginAndPassword() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestHandler_PostTextData(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		data        string
		expectedErr error
	}{
		{
			name:        "successful_save",
			token:       "test_token",
			data:        "test",
			expectedErr: nil,
		},
		{
			name:        "failed_null_resource",
			token:       "token",
			data:        "",
			expectedErr: ErrNotEmpty,
		},
		{
			name:        "failed_save",
			token:       "test_token",
			data:        "test",
			expectedErr: errSave,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			testServiceCredentials := NewMockserviceCredentials(ctrl)
			testServiceCards := NewMockserviceCards(ctrl)
			testServiceBinaryData := NewMockserviceBinaryData(ctrl)
			testServiceTextData := NewMockserviceTextData(ctrl)
			testServiceTextData.EXPECT().SaveTextData(ctx, tt.token, tt.data).Return(tt.expectedErr).AnyTimes()
			handler := NewHandlers(log, testServiceCredentials, testServiceTextData, testServiceBinaryData, testServiceCards)

			err := handler.PostTextData(ctx, tt.token, tt.data)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("PostLoginAndPassword() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestHandler_PostBinaryData(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		data        string
		expectedErr error
	}{
		{
			name:        "successful_save",
			token:       "test_token",
			data:        "test",
			expectedErr: nil,
		},
		{
			name:        "failed_null_resource",
			token:       "token",
			data:        "",
			expectedErr: ErrNotEmpty,
		},
		{
			name:        "failed_save",
			token:       "test_token",
			data:        "test",
			expectedErr: errSave,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			testServiceCredentials := NewMockserviceCredentials(ctrl)
			testServiceCards := NewMockserviceCards(ctrl)
			testServiceBinaryData := NewMockserviceBinaryData(ctrl)
			testServiceTextData := NewMockserviceTextData(ctrl)
			testServiceBinaryData.EXPECT().SaveBinaryData(ctx, tt.token, tt.data).Return(tt.expectedErr).AnyTimes()
			handler := NewHandlers(log, testServiceCredentials, testServiceTextData, testServiceBinaryData, testServiceCards)

			err := handler.PostBinaryData(ctx, tt.token, tt.data)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("PostLoginAndPassword() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestHandler_PostCards(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		data        string
		expectedErr error
	}{
		{
			name:        "successful_save",
			token:       "test_token",
			data:        "test",
			expectedErr: nil,
		},
		{
			name:        "failed_null_resource",
			token:       "token",
			data:        "",
			expectedErr: ErrNotEmpty,
		},
		{
			name:        "failed_save",
			token:       "test_token",
			data:        "test",
			expectedErr: errSave,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			testServiceCredentials := NewMockserviceCredentials(ctrl)
			testServiceCards := NewMockserviceCards(ctrl)
			testServiceBinaryData := NewMockserviceBinaryData(ctrl)
			testServiceTextData := NewMockserviceTextData(ctrl)
			testServiceCards.EXPECT().SaveCards(ctx, tt.token, tt.data).Return(tt.expectedErr).AnyTimes()
			handler := NewHandlers(log, testServiceCredentials, testServiceTextData, testServiceBinaryData, testServiceCards)

			err := handler.PostCards(ctx, tt.token, tt.data)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("PostLoginAndPassword() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}
