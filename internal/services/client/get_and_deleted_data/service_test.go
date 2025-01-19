package get_and_deleted_data

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
	errGetAll    = errors.New("failed to get all data")
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	db := NewMockstorage(ctrl)
	newService := NewService(log, db)

	if newService == nil {
		t.Errorf("Service is not nil")
	}
}

func TestGetAll_GetAllData(t *testing.T) {
	tests := []struct {
		name                 string
		expectedGetUserIDErr error
		expectedGetAllErr    error
	}{
		{
			name: "successful_get",
		},
		{
			name:                 "failed_get_user_id",
			expectedGetUserIDErr: errGetUserID,
		},
		{
			name:              "failed_get_all",
			expectedGetAllErr: errGetAll,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))

			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			db := NewMockstorage(ctrl)
			testSerivice := NewService(log, db)

			if test.expectedGetUserIDErr != nil {
				db.EXPECT().GetUserIDWithToken(ctx, gomock.Any()).Return(0, test.expectedGetUserIDErr)
			} else {
				db.EXPECT().GetUserIDWithToken(ctx, gomock.Any()).Return(1, nil)

				if test.expectedGetAllErr != nil {
					db.EXPECT().GetAll(ctx, 1, gomock.Any()).Return(nil, test.expectedGetAllErr)
				} else {
					db.EXPECT().GetAll(ctx, 1, gomock.Any()).Return(nil, nil)
				}
			}

			_, err := testSerivice.GetAllData(ctx, "test_token", "test_table")
			if !errors.Is(err, test.expectedGetAllErr) && !errors.Is(err, test.expectedGetUserIDErr) {
				t.Errorf("Service.GetAllData() error = %v, wantErr %v or %v", err, test.expectedGetAllErr, test.expectedGetUserIDErr)
			}

		})
	}
}

func TestGetAll_DeletedData(t *testing.T) {
	tests := []struct {
		name                 string
		token                string
		tableName            string
		id                   int
		expectedGetUserIDErr error
		expectedDelErr       error
		expectedSaveSyncErr  error
	}{
		{
			name:      "successful_del",
			token:     "test_token",
			tableName: "test_table",
			id:        1,
		},
		{
			name:                 "failed_get_user_id",
			token:                "test_token",
			tableName:            "test_table",
			id:                   1,
			expectedGetUserIDErr: errGetUserID,
		},
		{
			name:           "failed_del",
			token:          "test_token",
			tableName:      "test_table",
			id:             1,
			expectedDelErr: errors.New("failed to delete data from database"),
		},
		{
			name:                "failed_save_sync",
			token:               "test_token",
			tableName:           "test_table",
			id:                  1,
			expectedSaveSyncErr: errors.New("failed to save sync"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))

			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			db := NewMockstorage(ctrl)
			testSerivice := NewService(log, db)

			if tt.expectedGetUserIDErr != nil {
				db.EXPECT().GetUserIDWithToken(ctx, tt.token).Return(0, tt.expectedGetUserIDErr)
			} else {
				db.EXPECT().GetUserIDWithToken(ctx, tt.token).Return(1, nil)

				if tt.expectedDelErr != nil {
					db.EXPECT().Deleted(ctx, tt.tableName, tt.id).Return(tt.expectedDelErr)
				} else {
					db.EXPECT().Deleted(ctx, tt.tableName, tt.id).Return(nil)
					if tt.expectedSaveSyncErr != nil {
						db.EXPECT().SaveSync(ctx, tt.tableName, 1, tt.id, "deleted").Return(tt.expectedSaveSyncErr)
					} else {
						db.EXPECT().SaveSync(ctx, tt.tableName, 1, tt.id, "deleted").Return(nil)
					}
				}
			}

			err := testSerivice.DeletedData(ctx, tt.token, tt.tableName, tt.id)
			if !errors.Is(err, tt.expectedDelErr) && !errors.Is(err, tt.expectedGetUserIDErr) && !errors.Is(err, tt.expectedSaveSyncErr) {
				t.Errorf("Service.DeletedData() error = %v, wantErr %v or %v or %v", err, tt.expectedDelErr, tt.expectedGetUserIDErr, tt.expectedSaveSyncErr)
			}
		})
	}
}
