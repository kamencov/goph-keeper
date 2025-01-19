package auth

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
	"time"
)

var (
	errValidate  = errors.New("failed to get user_id")
	errSaveToken = errors.New("failed to generate token")
)

func TestServiceAuth_Auth(t *testing.T) {
	tests := []struct {
		name             string
		login            string
		password         string
		expectedCheckErr bool
		expectedSaveErr  error
		expectedGenErr   error
		expectedErr      error
	}{
		{
			name:             "successful_auth",
			login:            "test",
			password:         "test",
			expectedCheckErr: true,
			expectedSaveErr:  nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_check",
			login:            "test",
			password:         "test",
			expectedCheckErr: false,
			expectedErr:      ErrNotFoundLogin,
		},
		{
			name:             "failed_gen_token",
			login:            "test",
			password:         "test",
			expectedCheckErr: true,
			expectedGenErr:   errSaveToken,
			expectedSaveErr:  errSaveToken,
			expectedErr:      errSaveToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			db := NewMockstorageAuth(ctrl)

			newService := NewServiceAuth([]byte("salt"), []byte("salt"), time.Hour, log, db)
			if tt.expectedCheckErr == false {
				db.EXPECT().CheckPassword(tt.login).Return("", false).Times(1)
			} else {
				passwordHash := newService.hashPassword(tt.password)
				db.EXPECT().CheckPassword(tt.login).Return(passwordHash, true).Times(1)

				if tt.expectedGenErr != nil {
					db.EXPECT().SaveTableUserAndUpdateToken(tt.login, gomock.Any()).Return(tt.expectedSaveErr).Times(1)
				} else {
					db.EXPECT().SaveTableUserAndUpdateToken(tt.login, gomock.Any()).Return(nil).Times(1)
				}
			}

			_, err := newService.Auth(tt.login, tt.password)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestServiceAuth_ValidateToken(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		expectedGetErr error
		expectedErr    error
	}{
		{
			name:  "successful_validate",
			token: "test_token",
		},
		{
			name:           "failed_get_user_id",
			token:          "test_token",
			expectedGetErr: errValidate,
			expectedErr:    errValidate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			db := NewMockstorageAuth(ctrl)
			ctx := context.Background()
			newService := NewServiceAuth([]byte("salt"), []byte("salt"), time.Hour, log, db)
			if tt.expectedGetErr != nil {
				db.EXPECT().GetUserIDByToken(ctx, tt.token).Return(0, tt.expectedGetErr).Times(1)
			} else {
				db.EXPECT().GetUserIDByToken(ctx, tt.token).Return(1, nil).Times(1)
			}

			_, err := newService.ValidateToken(ctx, tt.token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestServiceAuth_CreatTokenForUser(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		expectedErr error
	}{
		{
			name:        "successful_create_token",
			userID:      "1",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			db := NewMockstorageAuth(ctrl)

			newService := NewServiceAuth([]byte("salt"), []byte("salt"), time.Hour, log, db)

			_, err := newService.CreatTokenForUser(tt.userID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
