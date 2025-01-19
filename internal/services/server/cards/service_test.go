package cards

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

func TestNewServiceCards(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug}))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockstorageCards(ctrl)
	newService := NewServiceCards(log, db)

	if newService == nil {
		t.Errorf("Service is not nil")
	}
}

func TestServiceCards_SyncSaveCards(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		accessToken     string
		cards           string
		expectedGetErr  error
		expectedSaveErr error
	}{
		{
			name:        "successful_save",
			userID:      1,
			accessToken: "test_token",
			cards:       "test cards",
		},
		{
			name:           "failed_get_user_id",
			userID:         1,
			accessToken:    "test_token",
			cards:          "test cards",
			expectedGetErr: errors.New("failed to get user_id"),
		},
		{
			name:            "failed_save",
			userID:          1,
			accessToken:     "test_token",
			cards:           "test cards",
			expectedSaveErr: errors.New("failed to handlers data"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			db := NewMockstorageCards(ctrl)
			if tt.expectedGetErr != nil {
				db.EXPECT().GetUserIDByToken(gomock.Any(), gomock.Any()).Return(0, tt.expectedGetErr).Times(1)
			} else {
				db.EXPECT().GetUserIDByToken(gomock.Any(), gomock.Any()).Return(tt.userID, nil).Times(1)
			}

			if tt.expectedSaveErr != nil {
				db.EXPECT().SaveCards(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSaveErr).Times(1)
			} else if tt.expectedGetErr == nil {
				// ServerSaveLoginAndPasswordInCredentials вызывается только если нет ошибки в GetUserIDByToken
				db.EXPECT().SaveCards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}

			newService := NewServiceCards(log, db)
			err := newService.SyncSaveCards(ctx, tt.accessToken, tt.cards)
			if !errors.Is(err, tt.expectedSaveErr) && !errors.Is(err, tt.expectedGetErr) {
				t.Errorf("ServiceCards.SyncSaveCards() error = %v, wantErr %v", err, tt.expectedSaveErr)
			}
		})
	}
}

func TestService_SyncDelBinary(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		accessToken    string
		cards          string
		expectedGetErr error
		expectedDelErr error
	}{
		{
			name:        "successful_save",
			userID:      1,
			accessToken: "test_token",
			cards:       "test cards",
		},
		{
			name:           "failed_get_user_id",
			userID:         1,
			accessToken:    "test_token",
			cards:          "test cards",
			expectedGetErr: errors.New("failed to get user_id"),
		},
		{
			name:           "failed_save",
			userID:         1,
			accessToken:    "test_token",
			cards:          "test cards",
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

			db := NewMockstorageCards(ctrl)

			// Настройка моков в зависимости от ожидаемого результата
			if tt.expectedGetErr != nil {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(0, tt.expectedGetErr).Times(1)
			} else {
				db.EXPECT().GetUserIDByToken(ctx, tt.accessToken).Return(tt.userID, nil).Times(1)
			}

			if tt.expectedDelErr != nil {
				db.EXPECT().DeletedCards(ctx, tt.userID, tt.cards).Return(tt.expectedDelErr).Times(1)
			} else if tt.expectedGetErr == nil {
				// DeletedText вызывается только если нет ошибки в GetUserIDByToken
				db.EXPECT().DeletedCards(ctx, tt.userID, tt.cards).Return(nil).Times(1)
			}

			newService := NewServiceCards(log, db)

			err := newService.SyncDelBinary(ctx, tt.accessToken, tt.cards)
			if !errors.Is(err, tt.expectedGetErr) && !errors.Is(err, tt.expectedDelErr) {
				t.Errorf("expected error %v or %v, got %v", tt.expectedGetErr, tt.expectedDelErr, err)
			}
		})
	}
}
