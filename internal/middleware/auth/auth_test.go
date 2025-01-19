package auth

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"os"
	"testing"
)

func TestNewMiddleware(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serviceMock := NewMockservice(ctrl)

	middleware := NewMiddleware(log, serviceMock)

	if middleware == nil {
		t.Errorf("Middleware is nil")
	}
}

func TestMiddleware_UnaryInterceptor(t *testing.T) {
	tests := []struct {
		name                  string
		method                string
		userID                int
		ctx                   context.Context
		expectedVerifyUserErr error
		expectedCode          codes.Code
	}{
		{
			name:                  "successful_unary_auth_token_true",
			method:                "/goph_keeper_v1.SyncFromClient/SyncFromClientCredentials",
			userID:                1,
			ctx:                   context.Background(),
			expectedVerifyUserErr: nil,
			expectedCode:          codes.OK,
		},
		{
			name:                  "failed_invalid_token",
			method:                "/goph_keeper_v1.SyncFromClient/SyncFromClientCredentials",
			userID:                1,
			ctx:                   context.Background(),
			expectedVerifyUserErr: errors.New("invalid token"),
			expectedCode:          codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serviceMock := NewMockservice(ctrl)
			serviceMock.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Return(tt.userID, tt.expectedVerifyUserErr).AnyTimes()
			req := "test-request"
			info := &grpc.UnaryServerInfo{
				FullMethod: tt.method,
			}
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				userID := ctx.Value(UserIDContextKey)
				return map[string]interface{}{"userID": userID}, nil
			}
			middleware := NewMiddleware(log, serviceMock)

			_, err := middleware.UnaryInterceptor(tt.ctx, req, info, handler)

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
