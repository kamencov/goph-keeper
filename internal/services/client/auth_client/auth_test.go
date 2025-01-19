package auth_client

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

var (
	errGetUserID   = errors.New("failed to get user id")
	errSave        = errors.New("failed to save token")
	errUpdate      = errors.New("failed to update token")
	errGetPassword = errors.New("failed to get user password")
	errGetToken    = errors.New("failed to get user token")
)

func TestNewService(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug}))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockstorage(ctrl)
	err := NewService(log, db)

	if err == nil {
		t.Errorf("Service is not nil")
	}

}

func TestService_SaveTokenInBase(t *testing.T) {
	tests := []struct {
		name                 string
		login                string
		password             string
		token                string
		expectedGetUserIDErr error
		expectedSaveErr      error
		expectedUpdateErr    error
	}{
		{
			name:     "successful_save",
			login:    "test_login",
			password: "test_password",
			token:    "test_token",
		},
		{
			name:                 "failed_get_user_id",
			login:                "test_login",
			password:             "test_password",
			token:                "test_token",
			expectedGetUserIDErr: errGetUserID,
		},
		{
			name:            "failed_save",
			login:           "test_login",
			password:        "test_password",
			token:           "test_token",
			expectedSaveErr: errSave,
		},
		{
			name:              "failed_update",
			login:             "test_login",
			password:          "test_password",
			token:             "test_token",
			expectedUpdateErr: errUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.Background()
			db := NewMockstorage(ctrl)
			testService := NewService(log, db)
			if tt.expectedGetUserIDErr != nil {
				db.EXPECT().GetUserIDWithLogin(ctx, tt.login).Return(0, tt.expectedGetUserIDErr).AnyTimes()
				if tt.expectedSaveErr != nil {
					db.EXPECT().SaveLoginAndToken(ctx, tt.login, tt.password, tt.token).Return(tt.expectedSaveErr).AnyTimes()
				} else {
					db.EXPECT().SaveLoginAndToken(ctx, tt.login, tt.password, tt.token).Return(nil).AnyTimes()
				}
			} else {
				db.EXPECT().GetUserIDWithLogin(ctx, tt.login).Return(1, nil).AnyTimes()

				if tt.expectedUpdateErr != nil {
					db.EXPECT().UpdateLoginAndToken(ctx, 1, tt.token).Return(tt.expectedUpdateErr).AnyTimes()
				} else {
					db.EXPECT().UpdateLoginAndToken(ctx, 1, tt.token).Return(nil).AnyTimes()
				}
			}

			err := testService.SaveTokenInBase(ctx, tt.login, tt.password, tt.token)

			if !errors.Is(err, tt.expectedGetUserIDErr) && !errors.Is(err, tt.expectedSaveErr) && !errors.Is(err, tt.expectedUpdateErr) {
				t.Errorf("Service is not nil")
			}
		})
	}
}

func TestService_CheckUser(t *testing.T) {
	tests := []struct {
		name                       string
		login                      string
		password                   string
		expectedGetUserPasswordErr error
		expectedGetUserTokenErr    error
	}{
		{
			name:     "successful_check",
			login:    "test_login",
			password: "test_password",
		},
		{
			name:                       "failed_get_user_password",
			login:                      "test_login",
			password:                   "test_password",
			expectedGetUserPasswordErr: errGetPassword,
		},
		{
			name:                    "failed_get_user_token",
			login:                   "test_login",
			password:                "test_password",
			expectedGetUserTokenErr: errGetToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.Background()
			db := NewMockstorage(ctrl)
			testService := NewService(log, db)

			if tt.expectedGetUserPasswordErr != nil {
				db.EXPECT().GetUserPassword(ctx, tt.login).Return("", tt.expectedGetUserPasswordErr).AnyTimes()
			} else {
				db.EXPECT().GetUserPassword(ctx, tt.login).Return(tt.password, nil).AnyTimes()
				if tt.expectedGetUserTokenErr != nil {
					db.EXPECT().GetUserToken(ctx, tt.login).Return("", tt.expectedGetUserTokenErr).AnyTimes()
				} else {
					db.EXPECT().GetUserToken(ctx, tt.login).Return(tt.password, nil).AnyTimes()
				}

				_, err := testService.CheckUser(ctx, tt.login, tt.password)
				if !errors.Is(err, tt.expectedGetUserPasswordErr) && !errors.Is(err, tt.expectedGetUserTokenErr) {
					t.Errorf("Service is not nil")
				}
			}
		})
	}
}
