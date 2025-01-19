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

func TestServiceAuth_RegisterUser(t *testing.T) {
	tests := []struct {
		name             string
		login            string
		password         string
		expectedCheckErr error
		expectedSaveErr  error
		expectedGetErr   error
	}{
		{
			name:     "successful_register",
			login:    "test",
			password: "test",
		},
		{
			name:             "failed_check",
			login:            "test",
			password:         "test",
			expectedCheckErr: errors.New("failed to check user"),
		},
		{
			name:            "failed_save",
			login:           "test",
			password:        "test",
			expectedSaveErr: errors.New("failed to save user"),
		},
		{
			name:           "failed_get",
			login:          "test",
			password:       "test",
			expectedGetErr: errors.New("failed to get user_id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))

			db := NewMockstorageAuth(ctrl)

			// Настройка моков в зависимости от сценария теста
			if tt.expectedCheckErr != nil {
				db.EXPECT().CheckUser(gomock.Any(), tt.login).Return(tt.expectedCheckErr).Times(1)
			} else {
				db.EXPECT().CheckUser(gomock.Any(), tt.login).Return(nil).Times(1)

				if tt.expectedSaveErr != nil {
					db.EXPECT().SaveUser(gomock.Any(), tt.login, gomock.Any()).Return(tt.expectedSaveErr).Times(1)
				} else {
					db.EXPECT().SaveUser(gomock.Any(), tt.login, gomock.Any()).Return(nil).Times(1)

					if tt.expectedGetErr != nil {
						db.EXPECT().GetUserIDByLogin(gomock.Any(), tt.login).Return(1, tt.expectedGetErr).Times(1)
					} else {
						db.EXPECT().GetUserIDByLogin(gomock.Any(), tt.login).Return(1, nil).Times(1)
					}
				}
			}

			testService := NewServiceAuth([]byte("salt"), []byte("salt"), time.Hour, log, db)

			_, err := testService.RegisterUser(context.Background(), tt.login, tt.password)

			if !errors.Is(err, tt.expectedCheckErr) && !errors.Is(err, tt.expectedSaveErr) && !errors.Is(err, tt.expectedGetErr) {
				t.Errorf("ServiceAuth.RegisterUser() error = %v, wantErr %v", err, tt.expectedCheckErr)
			}

		})
	}
}
