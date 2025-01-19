package sync

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1_pd "goph-keeper/internal/proto/v1"
	"log/slog"
	"os"
	"testing"
)

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	serviceCredMock := NewMockserviceCred(ctrl)
	serviceTextMock := NewMockserviceTextData(ctrl)
	serviceBinMock := NewMockserviceBinaryData(ctrl)
	serviceCardMock := NewMockserviceCards(ctrl)

	handler := NewHandler(log, serviceCredMock, serviceTextMock, serviceBinMock, serviceCardMock)

	if handler == nil {
		t.Errorf("Handler is nil")
	}
}

func TestHandler_SyncFromClientCredentials(t *testing.T) {
	tests := []struct {
		name            string
		in              *v1_pd.SyncFromClientCredentialsRequest
		expectedSaveErr error
		expectedDelErr  error
		expectedCode    codes.Code
	}{
		{
			name: "successful_sync",
			in: &v1_pd.SyncFromClientCredentialsRequest{
				Task: []*v1_pd.Credentials{
					{
						Id:          1,
						IdUser:      1,
						Resource:    "test_resource",
						Login:       "test_login",
						Password:    "test_password",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "save",
						AccessToken: "test_token",
					},
					{
						Id:          2,
						IdUser:      1,
						Resource:    "test_resource",
						Login:       "test_login",
						Password:    "test_password",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "deleted",
						AccessToken: "test_token",
					},
				},
			},
			expectedSaveErr: nil,
			expectedDelErr:  nil,
			expectedCode:    codes.OK,
		},
		{
			name: "failed_save_credentials",
			in: &v1_pd.SyncFromClientCredentialsRequest{
				Task: []*v1_pd.Credentials{
					{
						Id:          1,
						IdUser:      1,
						Resource:    "test_resource",
						Login:       "test_login",
						Password:    "test_password",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "save",
						AccessToken: "test_token",
					},
				},
			},
			expectedSaveErr: errors.New("failed to handlers data"),
			expectedCode:    codes.Internal,
		},
		{
			name: "failed_deleted_credentials",
			in: &v1_pd.SyncFromClientCredentialsRequest{
				Task: []*v1_pd.Credentials{
					{
						Id:          1,
						IdUser:      1,
						Resource:    "test_resource",
						Login:       "test_login",
						Password:    "test_password",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "deleted",
						AccessToken: "test_token",
					},
				},
			},
			expectedDelErr: errors.New("failed to handlers data"),
			expectedCode:   codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			serviceCredMock := NewMockserviceCred(ctrl)
			serviceTextMock := NewMockserviceTextData(ctrl)
			serviceBinMock := NewMockserviceBinaryData(ctrl)
			serviceCardMock := NewMockserviceCards(ctrl)

			serviceCredMock.EXPECT().SyncSaveCredentials(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedSaveErr).AnyTimes()
			serviceCredMock.EXPECT().SyncDelCredentials(ctx, gomock.Any(), gomock.Any()).Return(tt.expectedDelErr).AnyTimes()
			handler := NewHandler(log, serviceCredMock, serviceTextMock, serviceBinMock, serviceCardMock)
			_, err := handler.SyncFromClientCredentials(ctx, tt.in)
			if err != nil {
				code, ok := status.FromError(err)
				if !ok {
					t.Errorf("unexpected error type: %v", err)
				}
				if code.Code() != tt.expectedCode {
					t.Errorf("unexpected error code: got %v, want %v", code.Code(), tt.expectedCode)
				}

			} else if tt.expectedCode != codes.OK {
				t.Errorf("expected error code %v, got none", tt.expectedCode)
			}
		})
	}
}

func TestHandler_SyncFromClientTextData(t *testing.T) {
	tests := []struct {
		name            string
		in              *v1_pd.SyncFromClientTextDataRequest
		expectedSaveErr error
		expectedDelErr  error
		expectedCode    codes.Code
	}{
		{
			name: "successful_sync",
			in: &v1_pd.SyncFromClientTextDataRequest{
				Task: []*v1_pd.TextData{
					{
						Id:          1,
						IdUser:      1,
						Text:        "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "save",
						AccessToken: "test_token",
					},
					{
						Id:          2,
						IdUser:      1,
						Text:        "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "deleted",
						AccessToken: "test_token",
					},
				},
			},
			expectedSaveErr: nil,
			expectedDelErr:  nil,
			expectedCode:    codes.OK,
		},
		{
			name: "failed_save_credentials",
			in: &v1_pd.SyncFromClientTextDataRequest{
				Task: []*v1_pd.TextData{
					{
						Id:          1,
						IdUser:      1,
						Text:        "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "save",
						AccessToken: "test_token",
					},
				},
			},
			expectedSaveErr: errors.New("failed to handlers data"),
			expectedCode:    codes.Internal,
		},
		{
			name: "failed_deleted_credentials",
			in: &v1_pd.SyncFromClientTextDataRequest{
				Task: []*v1_pd.TextData{
					{
						Id:          1,
						IdUser:      1,
						Text:        "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "deleted",
						AccessToken: "test_token",
					},
				},
			},
			expectedDelErr: errors.New("failed to handlers data"),
			expectedCode:   codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			serviceCredMock := NewMockserviceCred(ctrl)
			serviceTextMock := NewMockserviceTextData(ctrl)
			serviceBinMock := NewMockserviceBinaryData(ctrl)
			serviceCardMock := NewMockserviceCards(ctrl)

			serviceTextMock.EXPECT().SyncSaveText(ctx, gomock.Any(), gomock.Any()).Return(tt.expectedSaveErr).AnyTimes()
			serviceTextMock.EXPECT().SyncDelText(ctx, gomock.Any(), gomock.Any()).Return(tt.expectedDelErr).AnyTimes()
			handler := NewHandler(log, serviceCredMock, serviceTextMock, serviceBinMock, serviceCardMock)
			_, err := handler.SyncFromClientTextData(ctx, tt.in)
			if err != nil {
				code, ok := status.FromError(err)
				if !ok {
					t.Errorf("unexpected error type: %v", err)
				}
				if code.Code() != tt.expectedCode {
					t.Errorf("unexpected error code: got %v, want %v", code.Code(), tt.expectedCode)
				}

			} else if tt.expectedCode != codes.OK {
				t.Errorf("expected error code %v, got none", tt.expectedCode)
			}
		})
	}
}

func TestHandler_SyncFromClientBinaryData(t *testing.T) {
	tests := []struct {
		name            string
		in              *v1_pd.SyncFromClientBinaryDataRequest
		expectedSaveErr error
		expectedDelErr  error
		expectedCode    codes.Code
	}{
		{
			name: "successful_sync",
			in: &v1_pd.SyncFromClientBinaryDataRequest{
				Task: []*v1_pd.BinaryData{
					{
						Id:          1,
						IdUser:      1,
						Binary:      "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "save",
						AccessToken: "test_token",
					},
					{
						Id:          2,
						IdUser:      1,
						Binary:      "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "deleted",
						AccessToken: "test_token",
					},
				},
			},
			expectedSaveErr: nil,
			expectedDelErr:  nil,
			expectedCode:    codes.OK,
		},
		{
			name: "failed_save_credentials",
			in: &v1_pd.SyncFromClientBinaryDataRequest{
				Task: []*v1_pd.BinaryData{
					{
						Id:          1,
						IdUser:      1,
						Binary:      "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "save",
						AccessToken: "test_token",
					},
				},
			},
			expectedSaveErr: errors.New("failed to handlers data"),
			expectedCode:    codes.Internal,
		},
		{
			name: "failed_deleted_credentials",
			in: &v1_pd.SyncFromClientBinaryDataRequest{
				Task: []*v1_pd.BinaryData{
					{
						Id:          1,
						IdUser:      1,
						Binary:      "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "deleted",
						AccessToken: "test_token",
					},
				},
			},
			expectedDelErr: errors.New("failed to handlers data"),
			expectedCode:   codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			serviceCredMock := NewMockserviceCred(ctrl)
			serviceTextMock := NewMockserviceTextData(ctrl)
			serviceBinMock := NewMockserviceBinaryData(ctrl)
			serviceCardMock := NewMockserviceCards(ctrl)

			serviceBinMock.EXPECT().SyncSaveBinary(ctx, gomock.Any(), gomock.Any()).Return(tt.expectedSaveErr).AnyTimes()
			serviceBinMock.EXPECT().SyncDelBinary(ctx, gomock.Any(), gomock.Any()).Return(tt.expectedDelErr).AnyTimes()
			handler := NewHandler(log, serviceCredMock, serviceTextMock, serviceBinMock, serviceCardMock)
			_, err := handler.SyncFromClientBinaryData(ctx, tt.in)
			if err != nil {
				code, ok := status.FromError(err)
				if !ok {
					t.Errorf("unexpected error type: %v", err)
				}
				if code.Code() != tt.expectedCode {
					t.Errorf("unexpected error code: got %v, want %v", code.Code(), tt.expectedCode)
				}

			} else if tt.expectedCode != codes.OK {
				t.Errorf("expected error code %v, got none", tt.expectedCode)
			}
		})
	}
}

func TestHandler_SyncFromClientCards(t *testing.T) {
	tests := []struct {
		name            string
		in              *v1_pd.SyncFromClientCardsRequest
		expectedSaveErr error
		expectedDelErr  error
		expectedCode    codes.Code
	}{
		{
			name: "successful_sync",
			in: &v1_pd.SyncFromClientCardsRequest{
				Task: []*v1_pd.Cards{
					{
						Id:          1,
						IdUser:      1,
						Cards:       "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "save",
						AccessToken: "test_token",
					},
					{
						Id:          2,
						IdUser:      1,
						Cards:       "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "deleted",
						AccessToken: "test_token",
					},
				},
			},
			expectedSaveErr: nil,
			expectedDelErr:  nil,
			expectedCode:    codes.OK,
		},
		{
			name: "failed_save_credentials",
			in: &v1_pd.SyncFromClientCardsRequest{
				Task: []*v1_pd.Cards{
					{
						Id:          1,
						IdUser:      1,
						Cards:       "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "save",
						AccessToken: "test_token",
					},
				},
			},
			expectedSaveErr: errors.New("failed to handlers data"),
			expectedCode:    codes.Internal,
		},
		{
			name: "failed_deleted_credentials",
			in: &v1_pd.SyncFromClientCardsRequest{
				Task: []*v1_pd.Cards{
					{
						Id:          1,
						IdUser:      1,
						Cards:       "test_text",
						UpdatedAt:   "2022-01-01T00:00:00Z",
						Action:      "deleted",
						AccessToken: "test_token",
					},
				},
			},
			expectedDelErr: errors.New("failed to handlers data"),
			expectedCode:   codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			serviceCredMock := NewMockserviceCred(ctrl)
			serviceTextMock := NewMockserviceTextData(ctrl)
			serviceBinMock := NewMockserviceBinaryData(ctrl)
			serviceCardMock := NewMockserviceCards(ctrl)

			serviceCardMock.EXPECT().SyncSaveCards(ctx, gomock.Any(), gomock.Any()).Return(tt.expectedSaveErr).AnyTimes()
			serviceCardMock.EXPECT().SyncDelBinary(ctx, gomock.Any(), gomock.Any()).Return(tt.expectedDelErr).AnyTimes()
			handler := NewHandler(log, serviceCredMock, serviceTextMock, serviceBinMock, serviceCardMock)
			_, err := handler.SyncFromClientCards(ctx, tt.in)
			if err != nil {
				code, ok := status.FromError(err)
				if !ok {
					t.Errorf("unexpected error type: %v", err)
				}
				if code.Code() != tt.expectedCode {
					t.Errorf("unexpected error code: got %v, want %v", code.Code(), tt.expectedCode)
				}

			} else if tt.expectedCode != codes.OK {
				t.Errorf("expected error code %v, got none", tt.expectedCode)
			}
		})
	}
}
