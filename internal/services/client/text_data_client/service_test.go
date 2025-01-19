package text_data_client

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
	errSaveText  = errors.New("failed to save text data in database")
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
	db := NewMockstorageTextDataClient(ctrl)
	err := NewService(log, db)

	if err == nil {
		t.Errorf("Service is not nil")
	}
}

func TestServiceClient_SaveTextData(t *testing.T) {
	tests := []struct {
		name                    string
		token                   string
		data                    string
		expectedGetUserIDErr    error
		expectedSaveTextDataErr error
		expectedGetIDTaskErr    error
		expectedSaveSyncErr     error
		expectedErr             error
	}{
		{
			name:        "successful_save",
			token:       "test_token",
			data:        "test_data",
			expectedErr: nil,
		},
		{
			name:                 "failed_get_user_id",
			token:                "test_token",
			data:                 "test_data",
			expectedGetUserIDErr: errGetUserID,
			expectedErr:          errGetUserID,
		},
		{
			name:                    "failed_save_text_data",
			token:                   "test_token",
			data:                    "test_data",
			expectedSaveTextDataErr: errSaveText,
			expectedErr:             errSaveText,
		},
		{
			name:                 "failed_get_id_task",
			token:                "test_token",
			data:                 "test_data",
			expectedGetIDTaskErr: errGetIDTask,
			expectedErr:          errGetIDTask,
		},
		{
			name:                "failed_save_sync",
			token:               "test_token",
			data:                "test_data",
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
			db := NewMockstorageTextDataClient(ctrl)

			ctx := context.Background()

			if tt.expectedGetUserIDErr != nil {
				db.EXPECT().GetUserIDWithToken(ctx, gomock.Any()).Return(0, tt.expectedGetUserIDErr).AnyTimes()
			} else {
				db.EXPECT().GetUserIDWithToken(ctx, gomock.Any()).Return(1, nil).AnyTimes()

				if tt.expectedSaveTextDataErr != nil {
					db.EXPECT().SaveTextDataInDatabase(ctx, 1, gomock.Any()).Return(tt.expectedSaveTextDataErr).AnyTimes()
				} else {
					db.EXPECT().SaveTextDataInDatabase(ctx, 1, gomock.Any()).Return(nil).AnyTimes()

					if tt.expectedGetIDTaskErr != nil {
						db.EXPECT().GetIDTaskText(ctx, gomock.Any(), 1, gomock.Any()).Return(0, tt.expectedGetIDTaskErr).AnyTimes()

					} else {
						db.EXPECT().GetIDTaskText(ctx, gomock.Any(), 1, gomock.Any()).Return(1, nil).AnyTimes()

						if tt.expectedSaveSyncErr != nil {
							db.EXPECT().SaveSync(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSaveSyncErr).AnyTimes()
						} else {
							db.EXPECT().SaveSync(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
						}
					}
				}
			}

			testService := NewService(log, db)

			err := testService.SaveTextData(ctx, tt.token, tt.data)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("SaveTextData() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}
