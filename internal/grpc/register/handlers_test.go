package register

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1_pd "goph-keeper/internal/proto/v1"
	"goph-keeper/internal/storage/postgresql"
	"log/slog"
	"os"
	"testing"
)

func TestNewHandlers(t *testing.T) {
	log := slog.Logger{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serviceMock := NewMockserviceRegister(ctrl)

	handler := NewHandlers(&log, serviceMock)

	if handler == nil {
		t.Errorf("serviceRegister is nil")
	}
}

func TestHandlers_Register(t *testing.T) {
	cases := []struct {
		name         string
		login        string
		password     string
		uid          int
		expectedErr  error
		expectedCode codes.Code
	}{
		{
			name:         "successful_register",
			login:        "test",
			password:     "test",
			uid:          1,
			expectedErr:  nil,
			expectedCode: codes.OK,
		},
		{
			name:         "empty_req",
			login:        "",
			password:     "",
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "user_already",
			login:        "test",
			password:     "test",
			uid:          1,
			expectedErr:  postgresql.ErrUserAlreadyExists,
			expectedCode: codes.AlreadyExists,
		},
		{
			name:         "failed_to_register",
			login:        "test",
			password:     "test",
			expectedErr:  sql.ErrNoRows,
			expectedCode: codes.Internal,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			ctx := context.Background()
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serviceMock := NewMockserviceRegister(ctrl)
			serviceMock.EXPECT().RegisterUser(ctx, gomock.Any(), gomock.Any()).
				Return(cc.uid, cc.expectedErr).AnyTimes()

			handler := NewHandlers(log, serviceMock)

			req := &v1_pd.RegisterRequest{
				Login:    cc.login,
				Password: cc.password}

			_, err := handler.Register(ctx, req)
			if err != nil {
				code, ok := status.FromError(err)
				if !ok {
					t.Errorf("unexpected error type: %v", err)
				}
				if code.Code() != cc.expectedCode {
					t.Errorf("unexpected error code: got %v, want %v", code.Code(), cc.expectedCode)
				}

			} else if cc.expectedCode != codes.OK {
				t.Errorf("expected error code %v, got none", cc.expectedCode)
			}
		})
	}
}
