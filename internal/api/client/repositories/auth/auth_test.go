package auth

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	v1_pd "goph-keeper/internal/proto/v1"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"
)

// MockRegisterServer - моковая реализация сервера Register
type MockRegisterServer struct {
	v1_pd.UnimplementedRegisterServer
}

// MockAuthServer - моковая реализация сервера Auth
type MockAuthServer struct {
	v1_pd.UnimplementedAuthServer
}

// Register - имитация регистрации пользователя
func (m *MockRegisterServer) Register(ctx context.Context, req *v1_pd.RegisterRequest) (*v1_pd.RegisterResponse, error) {
	if req.Login == "existing_user" {
		return nil, grpc.Errorf(409, "user already exists")
	}
	return &v1_pd.RegisterResponse{}, nil
}

// Auth - имитация авторизации пользователя
func (m *MockAuthServer) Auth(ctx context.Context, req *v1_pd.AuthRequest) (*v1_pd.AuthResponse, error) {
	if req.Login == "invalid_user" {
		return nil, grpc.Errorf(401, "invalid credentials")
	}
	return &v1_pd.AuthResponse{
		Token: "mocked_token",
	}, nil
}

var (
	errAuth = errors.New("auth error")
	errSave = errors.New("failed to save")
)

func TestNewHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	serviceMock := NewMockservice(ctrl)

	handler := NewHandlers(logs, serviceMock)

	if handler == nil {
		t.Errorf("Handler is nil")
	}
}

func TestHandlers_AuthUserOffLine(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		token       string
		password    string
		expectedErr error
	}{
		{
			name:        "successful_auth",
			login:       "test_login",
			token:       "test_token",
			password:    "test_password",
			expectedErr: nil,
		},
		{
			name:        "failed_auth",
			login:       "test_login",
			token:       "",
			password:    "test_password",
			expectedErr: errAuth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logs := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))

			ctx := context.Background()
			serviceMock := NewMockservice(ctrl)
			serviceMock.EXPECT().CheckUser(ctx, tt.login, tt.password).Return(tt.token, tt.expectedErr)

			handler := NewHandlers(logs, serviceMock)

			_, err := handler.AuthUserOffLine(ctx, tt.login, tt.password)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestHandlers_RegisterUser(t *testing.T) {
	// Настраиваем логгер
	logs := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	// Настройка буферного соединения для тестов
	const bufSize = 1024 * 1024
	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()

	// Регистрируем моковый сервер
	v1_pd.RegisterRegisterServer(server, &MockRegisterServer{})

	// Запускаем gRPC сервер в горутине
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	defer server.Stop()

	// Создаем соединение через буферный listener
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// Инициализируем объект Handlers
	h := &Handlers{
		log: logs,
	}

	tests := []struct {
		name      string
		login     string
		password  string
		expectErr bool
	}{
		{
			name:      "successful_auth",
			login:     "existing_user",
			password:  "password123",
			expectErr: false,
		},
		{
			name:      "failed_auth",
			login:     "non_existing_user",
			password:  "wrong_password",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := h.RegisterUser(ctx, conn, tt.login, tt.password)
			if (err == nil) != tt.expectErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestHandlers_AuthUser(t *testing.T) {
	// Настраиваем логгер
	logs := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	// Настройка буферного соединения для тестов
	const bufSize = 1024 * 1024
	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()

	// Регистрируем моковый сервер
	v1_pd.RegisterAuthServer(server, &MockAuthServer{})

	// Запускаем gRPC сервер в горутине
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	defer server.Stop()

	// Создаем соединение через буферный listener
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	tests := []struct {
		name            string
		login           string
		password        string
		expectedSaveErr error
		expectErr       bool
	}{
		{
			name:            "successful_auth",
			login:           "existing_user",
			password:        "password123",
			expectedSaveErr: nil,
			expectErr:       false,
		},
		{
			name:            "failed_auth_save",
			login:           "existing_user",
			password:        "password123",
			expectedSaveErr: errSave,
			expectErr:       true,
		},
		{
			name:            "failed_auth",
			login:           "invalid_user",
			password:        "wrong_password",
			expectedSaveErr: nil,
			expectErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// читаем из канала SignalAuth сигнал об авторизации
			go func() {
				<-SignalAuth

			}()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serviceMock := NewMockservice(ctrl)
			serviceMock.EXPECT().
				SaveTokenInBase(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tt.expectedSaveErr).
				AnyTimes()
			// Инициализируем объект Handlers
			h := &Handlers{
				log:     logs,
				service: serviceMock,
			}
			_, err := h.AuthUser(ctx, conn, tt.login, tt.password)

			if (err != nil) != tt.expectErr {
				t.Errorf("AuthUser() error = %v, wantErr %v", err, tt.expectErr)
			}

		})
	}
}
