package text_data

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
	db := NewMockstorageTextData(ctrl)
	newService := NewService(log, db)

	if newService == nil {
		t.Errorf("Service is not nil")
	}
}

func TestService_SyncSaveText(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		accessToken     string
		data            string
		expectedGetErr  error
		expectedSaveErr error
	}{
		{
			name:            "successful_save",
			userID:          1,
			accessToken:     "test_token",
			data:            "test data",
			expectedGetErr:  nil,
			expectedSaveErr: nil,
		},
		{
			name:        "failed_get_user_id",
			userID:      1,
			accessToken: "test_token",
			data:        "test data",

			expectedGetErr: errors.New("failed to get user_id"),
		},
		{
			name:            "failed_save",
			userID:          1,
			accessToken:     "test_token",
			data:            "test data",
			expectedGetErr:  nil,
			expectedSaveErr: errors.New("failed to save text data"),
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

			db := NewMockstorageTextData(ctrl)

			// Настройка моков в зависимости от ожидаемого результата
			if tt.expectedGetErr != nil {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(0, tt.expectedGetErr).Times(1)
			} else {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(tt.userID, nil).Times(1)
			}

			if tt.expectedSaveErr != nil {
				db.EXPECT().SaveTextDataPstgres(ctx, tt.userID, tt.data).Return(tt.expectedSaveErr).Times(1)
			} else if tt.expectedGetErr == nil {
				// SaveTextDataPstgres вызывается только если нет ошибки в GetUserIDByToken
				db.EXPECT().SaveTextDataPstgres(ctx, tt.userID, tt.data).Return(nil).Times(1)
			}

			newService := NewService(log, db)

			err := newService.SyncSaveText(ctx, tt.accessToken, tt.data)
			if !errors.Is(err, tt.expectedGetErr) && !errors.Is(err, tt.expectedSaveErr) {
				t.Errorf("expected error %v or %v, got %v", tt.expectedGetErr, tt.expectedSaveErr, err)
			}
		})
	}
}

func TestService_SyncDelText(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		accessToken    string
		data           string
		expectedGetErr error
		expectedDelErr error
	}{
		{
			name:           "successful_del",
			userID:         1,
			accessToken:    "test_token",
			data:           "test data",
			expectedGetErr: nil,
			expectedDelErr: nil,
		},
		{
			name:           "failed_get_user_id",
			userID:         1,
			accessToken:    "test_token",
			data:           "test data",
			expectedGetErr: errors.New("failed to get user_id"),
		},
		{
			name:           "failed_del",
			userID:         1,
			accessToken:    "test_token",
			data:           "test data",
			expectedGetErr: nil,
			expectedDelErr: errors.New("failed to deleted text"),
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

			db := NewMockstorageTextData(ctrl)

			// Настройка моков в зависимости от ожидаемого результата
			if tt.expectedGetErr != nil {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(0, tt.expectedGetErr).Times(1)
			} else {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(tt.userID, nil).Times(1)
			}

			if tt.expectedDelErr != nil {
				db.EXPECT().DeletedText(ctx, tt.userID, tt.data).Return(tt.expectedDelErr).Times(1)
			} else if tt.expectedGetErr == nil {
				// DeletedText вызывается только если нет ошибки в GetUserIDByToken
				db.EXPECT().DeletedText(ctx, tt.userID, tt.data).Return(nil).Times(1)
			}

			newService := NewService(log, db)

			err := newService.SyncDelText(ctx, tt.accessToken, tt.data)
			if !errors.Is(err, tt.expectedGetErr) && !errors.Is(err, tt.expectedDelErr) {
				t.Errorf("expected error %v or %v, got %v", tt.expectedGetErr, tt.expectedDelErr, err)
			}
		})
	}
}
