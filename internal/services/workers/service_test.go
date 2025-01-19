package workers

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

var (
	errGetAllErr   = errors.New("failed to get all data")
	errGetDataErr  = errors.New("failed to get data")
	errGetTokenErr = errors.New("failed to get token")
	errSyncErr     = errors.New("failed to sync")
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	db := NewMockstorage(ctrl)
	dbSync := NewMockrepoSync(ctrl)
	newService := NewService(log, db, dbSync)

	if newService == nil {
		t.Errorf("Service is not nil")
	}
}

func TestService_PushData(t *testing.T) {
	tests := []struct {
		name                 string
		syncsData            []*SyncModel
		expectedGetAlLErr    error
		LenData              int
		expectedGetDataErr   error
		expectedGetTokenErr  error
		expectedSyncErr      error
		expectedClearSyncErr error
		expectedErr          error
	}{
		{
			name: "test",
			syncsData: []*SyncModel{
				{
					UserID: 1,
					TaskID: 1,
					Action: "test",
				},
			},
			LenData: 1,
		},
		{
			name:              "failed_get_all",
			expectedGetAlLErr: errGetAllErr,
			expectedErr:       errGetAllErr,
		},
		{
			name:        "failed_len_data",
			LenData:     0,
			expectedErr: errLen,
		},
		{
			name: "failed_get_data_credentials",
			syncsData: []*SyncModel{
				{
					UserID:    1,
					TableName: "credentials",
				},
			},
			LenData:            1,
			expectedGetDataErr: errGetDataErr,
			expectedErr:        errGetDataErr,
		},
		{
			name: "failed_get_data_text_data",
			syncsData: []*SyncModel{
				{
					UserID:    1,
					TableName: "text_data",
				},
			},
			LenData:            1,
			expectedGetDataErr: errGetDataErr,
			expectedErr:        errGetDataErr,
		},
		{
			name: "failed_get_data_binary_data",
			syncsData: []*SyncModel{
				{
					UserID:    1,
					TableName: "binary_data",
				},
			},
			LenData:            1,
			expectedGetDataErr: errGetDataErr,
			expectedErr:        errGetDataErr,
		},
		{
			name: "failed_get_data_cards",
			syncsData: []*SyncModel{
				{
					UserID:    1,
					TableName: "cards",
				},
			},
			LenData:            1,
			expectedGetDataErr: errGetDataErr,
			expectedErr:        errGetDataErr,
		},
		{
			name: "failed_get_token_credentials",
			syncsData: []*SyncModel{
				{
					UserID:    1,
					TableName: "credentials",
				},
			},
			LenData:             1,
			expectedGetTokenErr: errGetTokenErr,
			expectedErr:         errGetTokenErr,
		},
		{
			name: "failed_get_token_text_data",
			syncsData: []*SyncModel{
				{
					UserID:    1,
					TableName: "text_data",
				},
			},
			LenData:             1,
			expectedGetTokenErr: errGetTokenErr,
			expectedErr:         errGetTokenErr,
		},
		{
			name: "failed_get_token_binary_data",
			syncsData: []*SyncModel{
				{
					UserID:    1,
					TableName: "binary_data",
				},
			},
			LenData:             1,
			expectedGetTokenErr: errGetTokenErr,
			expectedErr:         errGetTokenErr,
		},
		{
			name: "failed_get_token_cards",
			syncsData: []*SyncModel{
				{
					UserID:    1,
					TableName: "cards",
				},
			},
			LenData:             1,
			expectedGetTokenErr: errGetTokenErr,
			expectedErr:         errGetTokenErr,
		},
		{
			name: "failed_sync",
			syncsData: []*SyncModel{
				{
					UserID:    1,
					TableName: "credentials",
				},
			},
			LenData:              1,
			expectedClearSyncErr: errClearSyncErr,
			expectedErr:          errClearSyncErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))

			db := NewMockstorage(ctrl)
			dbSync := NewMockrepoSync(ctrl)
			newService := NewService(log, db, dbSync)

			if tt.expectedGetAlLErr != nil {
				db.EXPECT().GetAllSync().Return(tt.syncsData, tt.expectedGetAlLErr)
			} else {
				if tt.LenData == 0 {
					db.EXPECT().GetAllSync().Return([]*SyncModel{}, nil).AnyTimes()
				} else {
					db.EXPECT().GetAllSync().Return(tt.syncsData, nil).AnyTimes()
					db.EXPECT().GetDataCredentials(gomock.Any(), gomock.Any()).Return(&Credentials{
						UserID: 1}, tt.expectedGetDataErr).AnyTimes()
					db.EXPECT().GetDataTextData(gomock.Any(), gomock.Any()).Return(&TextData{ID: 1}, tt.expectedGetDataErr).AnyTimes()
					db.EXPECT().GetDataBinaryData(gomock.Any(), gomock.Any()).Return(&BinaryData{ID: 1}, tt.expectedGetDataErr).AnyTimes()
					db.EXPECT().GetDataCards(gomock.Any(), gomock.Any()).Return(&Cards{ID: 1}, tt.expectedGetDataErr).AnyTimes()
					db.EXPECT().GetTokenWithUserID(gomock.Any(), gomock.Any()).Return("token", tt.expectedGetTokenErr).AnyTimes()
					dbSync.EXPECT().SyncCredentials(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSyncErr).AnyTimes()
					dbSync.EXPECT().SyncTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSyncErr).AnyTimes()
					dbSync.EXPECT().SyncBinaryData(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSyncErr).AnyTimes()
					dbSync.EXPECT().SyncCards(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSyncErr).AnyTimes()
					db.EXPECT().ClearSyncCredentials().Return(tt.expectedClearSyncErr).AnyTimes()
					db.EXPECT().ClearSyncTextData().Return(tt.expectedClearSyncErr).AnyTimes()
					db.EXPECT().ClearSyncBinaryData().Return(tt.expectedClearSyncErr).AnyTimes()
					db.EXPECT().ClearSyncCards().Return(tt.expectedClearSyncErr).AnyTimes()
				}
			}
			ctx := context.Background()
			err := newService.PushData(ctx, nil)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Service is not nil")
			}
		})
	}
}
