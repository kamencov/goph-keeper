package auth

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

var errAuth = errors.New("auth error")

func TestNewHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	serviceMock := NewMockservice(ctrl)

	handler := NewHandlers(log, serviceMock)

	if handler == nil {
		t.Errorf("Handler is nil")
	}
}

func TestHandlers_AuthUserOffLine(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		token       string
		password    string
		expectedErr error
	}{
		{
			name:        "successful_auth",
			login:       "test_login",
			token:       "test_token",
			password:    "test_password",
			expectedErr: nil,
		},
		{
			name:        "failed_auth",
			login:       "test_login",
			token:       "",
			password:    "test_password",
			expectedErr: errAuth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			ctx := context.Background()
			serviceMock := NewMockservice(ctrl)
			serviceMock.EXPECT().CheckUser(ctx, tt.login, tt.password).Return(tt.token, tt.expectedErr)

			handler := NewHandlers(log, serviceMock)

			_, err := handler.AuthUserOffLine(ctx, tt.login, tt.password)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
