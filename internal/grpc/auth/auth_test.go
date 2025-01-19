package auth

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1_pd "goph-keeper/internal/proto/v1"
	"goph-keeper/internal/services/server/auth"
	"log/slog"
	"os"
	"testing"
)

func TestNewHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	serviceMock := NewMockserviceAuth(ctrl)

	handler := NewHandlers(log, serviceMock)

	if handler == nil {
		t.Errorf("Handler is nil")
	}
}

func TestHandlers_Auth(t *testing.T) {
	tests := []struct {
		name            string
		in              *v1_pd.AuthRequest
		token           auth.Tokens
		expectedAuthErr error
		expectedCode    codes.Code
	}{
		{
			name: "successful_auth",
			in: &v1_pd.AuthRequest{
				Login:    "test_login",
				Password: "test_password",
			},
			token: auth.Tokens{
				AccessToken: "test_token"},
			expectedCode: codes.OK,
		},
		{
			name: "in_null",
			in: &v1_pd.AuthRequest{
				Login:    "",
				Password: "",
			},
			token:        auth.Tokens{},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "failed_auth_service",
			in: &v1_pd.AuthRequest{
				Login:    "test_login",
				Password: "test_password",
			},
			token:           auth.Tokens{},
			expectedAuthErr: errors.New("failed to auth user"),
			expectedCode:    codes.Internal,
		},
		{
			name: "failed_not_found_login",
			in: &v1_pd.AuthRequest{
				Login:    "test_login",
				Password: "test_password",
			},
			token:           auth.Tokens{},
			expectedAuthErr: auth.ErrNotFoundLogin,
			expectedCode:    codes.NotFound,
		},
		{
			name: "failed_wrong_password",
			in: &v1_pd.AuthRequest{
				Login:    "test_login",
				Password: "test_password",
			},
			token:           auth.Tokens{},
			expectedAuthErr: auth.ErrWrongPassword,
			expectedCode:    codes.Unauthenticated,
		},
		{
			name: "failed_not_token",
			in: &v1_pd.AuthRequest{
				Login:    "test_login",
				Password: "test_password",
			},
			token: auth.Tokens{
				AccessToken: "",
			},
			expectedAuthErr: nil,
			expectedCode:    codes.Internal,
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

			serviceMock := NewMockserviceAuth(ctrl)
			serviceMock.EXPECT().Auth(tt.in.Login, tt.in.Password).Return(tt.token, tt.expectedAuthErr).AnyTimes()

			handler := NewHandlers(log, serviceMock)

			_, err := handler.Auth(context.Background(), tt.in)

			if err != nil {
				code, ok := status.FromError(err)
				if !ok {
					t.Errorf("unexpected error type: %v", err)
				}
				if code.Code() != tt.expectedCode {
					t.Errorf("unexpected error code: got %v, want %v", code.Code(), tt.expectedCode)
				}

			} else if tt.expectedCode != codes.OK {
				t.Errorf("expected error code %v, got none", tt.expectedCode)
			}
		})
	}
}
