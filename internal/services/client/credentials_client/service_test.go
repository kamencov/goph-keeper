package credentials_client

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

var (
	errGetUserID = errors.New("failed to get user id")
	errSaveText  = errors.New("failed to save credentials data in database")
	errGetIDTask = errors.New("failed to get id task")
	errSaveSync  = errors.New("failed to save sync")
)

func TestNewService(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockcredentialsClient(ctrl)
	err := NewService(log, db)

	if err == nil {
		t.Errorf("Service is not nil")
	}
}

func TestServiceClient_SaveTextData(t *testing.T) {
	tests := []struct {
		name                 string
		token                string
		resource             string
		login                string
		password             string
		expectedGetUserIDErr error
		expectedSaveCardsErr error
		expectedGetIDTaskErr error
		expectedSaveSyncErr  error
		expectedErr          error
	}{
		{
			name:        "successful_save",
			token:       "test_token",
			resource:    "test_resource",
			login:       "test_login",
			password:    "test_password",
			expectedErr: nil,
		},
		{
			name:                 "failed_get_user_id",
			token:                "test_token",
			resource:             "test_resource",
			login:                "test_login",
			password:             "test_password",
			expectedGetUserIDErr: errGetUserID,
			expectedErr:          errGetUserID,
		},
		{
			name:                 "failed_save_credentials",
			token:                "test_token",
			resource:             "test_resource",
			login:                "test_login",
			password:             "test_password",
			expectedSaveCardsErr: errSaveText,
			expectedErr:          errSaveText,
		},
		{
			name:                 "failed_get_id_task",
			token:                "test_token",
			resource:             "test_resource",
			login:                "test_login",
			password:             "test_password",
			expectedGetIDTaskErr: errGetIDTask,
			expectedErr:          errGetIDTask,
		},
		{
			name:                "failed_save_sync",
			token:               "test_token",
			resource:            "test_resource",
			login:               "test_login",
			password:            "test_password",
			expectedSaveSyncErr: errSaveSync,
			expectedErr:         errSaveSync,
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
			db := NewMockcredentialsClient(ctrl)

			ctx := context.Background()

			if tt.expectedGetUserIDErr != nil {
				db.EXPECT().GetUserIDWithToken(ctx, gomock.Any()).Return(0, tt.expectedGetUserIDErr).AnyTimes()
			} else {
				db.EXPECT().GetUserIDWithToken(ctx, gomock.Any()).Return(1, nil).AnyTimes()

				if tt.expectedSaveCardsErr != nil {
					db.EXPECT().SaveLoginAndPasswordInCredentials(ctx, 1, tt.resource, tt.login, tt.password).Return(tt.expectedSaveCardsErr).AnyTimes()
				} else {
					db.EXPECT().SaveLoginAndPasswordInCredentials(ctx, 1, tt.resource, tt.login, tt.password).Return(nil).AnyTimes()

					if tt.expectedGetIDTaskErr != nil {
						db.EXPECT().GetIDTaskCredentials(ctx, gomock.Any(), 1, tt.resource).Return(0, tt.expectedGetIDTaskErr).AnyTimes()

					} else {
						db.EXPECT().GetIDTaskCredentials(ctx, gomock.Any(), 1, tt.resource).Return(1, nil).AnyTimes()

						if tt.expectedSaveSyncErr != nil {
							db.EXPECT().SaveSync(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSaveSyncErr).AnyTimes()
						} else {
							db.EXPECT().SaveSync(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
						}
					}
				}
			}

			testService := NewService(log, db)

			err := testService.SaveLoginAndPassword(ctx, tt.token, tt.resource, tt.login, tt.password)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("SaveTextData() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}
