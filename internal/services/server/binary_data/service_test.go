package binary_data

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

func TestNewService(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug}))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockstorage(ctrl)
	newService := NewService(log, db)

	if newService == nil {
		t.Errorf("Service is not nil")
	}
}

func TestService_SyncSaveBinary(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		accessToken     string
		data            string
		expectedGetErr  error
		expectedSaveErr error
	}{
		{
			name:        "successful_save",
			userID:      1,
			accessToken: "test_token",
			data:        "test data",
		},
		{
			name:           "failed_get_user_id",
			userID:         1,
			accessToken:    "test_token",
			data:           "test data",
			expectedGetErr: errors.New("failed to get user_id"),
		},
		{
			name:            "failed_save",
			userID:          1,
			accessToken:     "test_token",
			data:            "test data",
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
			db := NewMockstorage(ctrl)
			// Настройка моков в зависимости от ожидаемого результата
			if tt.expectedGetErr != nil {
				db.EXPECT().GetUserIDByToken(gomock.Any(), gomock.Any()).Return(0, tt.expectedGetErr).Times(1)
			} else {
				db.EXPECT().GetUserIDByToken(gomock.Any(), gomock.Any()).Return(tt.userID, nil).Times(1)
			}

			if tt.expectedSaveErr != nil {
				db.EXPECT().SaveBinaryDataBinary(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSaveErr).Times(1)
			} else if tt.expectedGetErr == nil {
				// ServerSaveLoginAndPasswordInCredentials вызывается только если нет ошибки в GetUserIDByToken
				db.EXPECT().SaveBinaryDataBinary(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}
			newService := NewService(log, db)
			ctx := context.Background()
			err := newService.SyncSaveBinary(ctx, tt.accessToken, tt.data)
			if !errors.Is(err, tt.expectedGetErr) && !errors.Is(err, tt.expectedSaveErr) {
				t.Errorf("expected error %v or %v, got %v", tt.expectedGetErr, tt.expectedSaveErr, err)
			}
		})
	}
}

func TestService_SyncDelBinary(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		accessToken    string
		data           string
		expectedGetErr error
		expectedDelErr error
	}{
		{
			name:        "successful_save",
			userID:      1,
			accessToken: "test_token",
			data:        "test data",
		},
		{
			name:           "failed_get_user_id",
			userID:         1,
			accessToken:    "test_token",
			data:           "test data",
			expectedGetErr: errors.New("failed to get user_id"),
		},
		{
			name:           "failed_save",
			userID:         1,
			accessToken:    "test_token",
			data:           "test data",
			expectedDelErr: errors.New("failed to handlers data"),
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

			db := NewMockstorage(ctrl)

			// Настройка моков в зависимости от ожидаемого результата
			if tt.expectedGetErr != nil {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(0, tt.expectedGetErr).Times(1)
			} else {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(tt.userID, nil).Times(1)
			}

			if tt.expectedDelErr != nil {
				db.EXPECT().DeletedBinary(ctx, tt.userID, tt.data).Return(tt.expectedDelErr).Times(1)
			} else if tt.expectedGetErr == nil {
				// DeletedText вызывается только если нет ошибки в GetUserIDByToken
				db.EXPECT().DeletedBinary(ctx, tt.userID, tt.data).Return(nil).Times(1)
			}

			newService := NewService(log, db)

			err := newService.SyncDelBinary(ctx, tt.accessToken, tt.data)
			if !errors.Is(err, tt.expectedGetErr) && !errors.Is(err, tt.expectedDelErr) {
				t.Errorf("expected error %v or %v, got %v", tt.expectedGetErr, tt.expectedDelErr, err)
			}
		})
	}
}
