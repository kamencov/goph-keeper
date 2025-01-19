package credentials

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

func TestNewService(t *testing.T) {
	var log *slog.Logger

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockcredentials(ctrl)
	newService := NewService(log, db)

	if newService == nil {
		t.Errorf("Service is not nil")
	}
}

type resource struct {
	accessToken string
	res         string
	login       string
	password    string
}

func TestService_SyncSaveCredentials(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		resource        *resource
		expectedGetErr  error
		expectedSaveErr error
	}{
		{
			name:            "successful_save",
			userID:          1,
			resource:        &resource{accessToken: "test_token", res: "test_resource", login: "test_login", password: "test_password"},
			expectedGetErr:  nil,
			expectedSaveErr: nil,
		},
		{
			name:           "failed_get_user_id",
			userID:         1,
			resource:       &resource{accessToken: "test_token", res: "test_resource", login: "test_login", password: "test_password"},
			expectedGetErr: ErrNotFoundUser,
		},
		{
			name:            "failed_save",
			userID:          1,
			resource:        &resource{accessToken: "test_token", res: "test_resource", login: "test_login", password: "test_password"},
			expectedGetErr:  nil,
			expectedSaveErr: errors.New("failed to handlers data"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelDebug}))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			db := NewMockcredentials(ctrl)
			// Настройка моков в зависимости от ожидаемого результата
			if tt.expectedGetErr != nil {
				db.EXPECT().GetUserIDByToken(gomock.Any(), gomock.Any()).Return(0, tt.expectedGetErr).Times(1)
			} else {
				db.EXPECT().GetUserIDByToken(gomock.Any(), gomock.Any()).Return(tt.userID, nil).Times(1)
			}

			if tt.expectedSaveErr != nil {
				db.EXPECT().ServerSaveLoginAndPasswordInCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSaveErr).Times(1)
			} else if tt.expectedGetErr == nil {
				// ServerSaveLoginAndPasswordInCredentials вызывается только если нет ошибки в GetUserIDByToken
				db.EXPECT().ServerSaveLoginAndPasswordInCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}
			newService := NewService(log, db)
			ctx := context.Background()
			err := newService.SyncSaveCredentials(ctx, tt.resource.accessToken, tt.resource.res, tt.resource.login, tt.resource.password)
			if !errors.Is(err, tt.expectedGetErr) && !errors.Is(err, tt.expectedSaveErr) {
				t.Errorf("expected error %v or %v, got %v", tt.expectedGetErr, tt.expectedSaveErr, err)
			}
		})
	}
}

func TestService_SyncDelCredentials(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		accessToken    string
		resource       string
		expectedGetErr error
		expectedDelErr error
	}{
		{
			name:           "successful_del",
			userID:         1,
			accessToken:    "test_token",
			resource:       "test_resource",
			expectedGetErr: nil,
			expectedDelErr: nil,
		},
		{
			name:           "failed_get_user_id",
			userID:         1,
			accessToken:    "test_token",
			resource:       "test_resource",
			expectedGetErr: ErrNotFoundUser,
		},
		{
			name:           "failed_del",
			userID:         1,
			accessToken:    "test_token",
			resource:       "test_resource",
			expectedGetErr: nil,
			expectedDelErr: errors.New("failed to deleted credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelDebug}))
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := NewMockcredentials(ctrl)

			// Настройка моков в зависимости от ожидаемого результата
			if tt.expectedGetErr != nil {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(0, tt.expectedGetErr).Times(1)
			} else {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(tt.userID, nil).Times(1)
			}

			if tt.expectedDelErr != nil {
				db.EXPECT().DeletedCredentials(ctx, tt.userID, tt.resource).Return(tt.expectedDelErr).Times(1)
			} else if tt.expectedGetErr == nil {
				// DeletedText вызывается только если нет ошибки в GetUserIDByToken
				db.EXPECT().DeletedCredentials(ctx, tt.userID, tt.resource).Return(nil).Times(1)
			}

			newService := NewService(log, db)

			err := newService.SyncDelCredentials(ctx, tt.accessToken, tt.resource)
			if !errors.Is(err, tt.expectedGetErr) && !errors.Is(err, tt.expectedDelErr) {
				t.Errorf("expected error %v or %v, got %v", tt.expectedGetErr, tt.expectedDelErr, err)
			}
		})
	}
}
