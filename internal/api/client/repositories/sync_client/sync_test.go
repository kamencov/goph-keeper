package sync_client

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	v1_pd "goph-keeper/internal/proto/v1"
	"goph-keeper/internal/services/workers"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"
)

var (
	errNoTask = errors.New("no task")
)

// MockSyncCredentialsServer - моковая реализация сервера SyncCredentials.
type MockSyncCredentialsServer struct {
	v1_pd.UnimplementedSyncFromClientServer
}

// SyncFromClientCredentials синхронизирует учётные данные.
func (m *MockSyncCredentialsServer) SyncFromClientCredentials(ctx context.Context, in *v1_pd.SyncFromClientCredentialsRequest) (*v1_pd.Empty, error) {
	task := in.Task[0]
	if task.Login == "invalid" {
		return nil, grpc.Errorf(401, "invalid credentials")
	}
	return &v1_pd.Empty{Message: "completed"}, nil
}

// SyncFromClientTextData синхронизирует учётные данные.
func (m *MockSyncCredentialsServer) SyncFromClientTextData(ctx context.Context, in *v1_pd.SyncFromClientTextDataRequest) (*v1_pd.Empty, error) {
	task := in.Task[0]
	if task.Text == "invalid" {
		return nil, grpc.Errorf(401, "invalid text data")
	}
	return &v1_pd.Empty{Message: "completed"}, nil
}

// SyncFromClientBinaryData синхронизирует учётные данные.
func (m *MockSyncCredentialsServer) SyncFromClientBinaryData(ctx context.Context, in *v1_pd.SyncFromClientBinaryDataRequest) (*v1_pd.Empty, error) {
	task := in.Task[0]
	if task.Binary == "invalid" {
		return nil, grpc.Errorf(401, "invalid text data")
	}
	return &v1_pd.Empty{Message: "completed"}, nil
}

// SyncFromClientCards синхронизирует учётные данные.
func (m *MockSyncCredentialsServer) SyncFromClientCards(ctx context.Context, in *v1_pd.SyncFromClientCardsRequest) (*v1_pd.Empty, error) {
	task := in.Task[0]
	if task.Cards == "invalid" {
		return nil, grpc.Errorf(401, "invalid text data")
	}
	return &v1_pd.Empty{Message: "completed"}, nil
}

func TestNewHandlers(t *testing.T) {
	logs := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	handler := NewHandlers(logs)

	if handler == nil {
		t.Errorf("Handler is nil")
	}
}

func TestHandler_SyncCredentials(t *testing.T) {
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
	v1_pd.RegisterSyncFromClientServer(server, &MockSyncCredentialsServer{})

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
		name         string
		data         []*workers.Credentials
		expectedErr  error
		expectedBool bool
	}{
		{
			name:         "successful_sync",
			data:         []*workers.Credentials{{ID: 1, UserID: 1, Resource: "test", Login: "test", Password: "test", UpdatedAt: time.Now(), Action: "test", AccessToken: "test"}},
			expectedBool: true,
		},
		{
			name:         "data_is_nil",
			expectedBool: false,
		},
		{
			name:         "failed_sync",
			data:         []*workers.Credentials{{ID: 1, UserID: 1, Resource: "test", Login: "invalid", Password: "test", UpdatedAt: time.Now(), Action: "test", AccessToken: "test"}},
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandlers(logs)
			err = handler.SyncCredentials(ctx, conn, tt.data)
			if (err != nil) == tt.expectedBool {
				t.Errorf("SyncCredentials() error = %v, wantErr %v", err, tt.expectedErr)

			}
		})
	}
}

func TestHandler_SyncTextData(t *testing.T) {
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
	v1_pd.RegisterSyncFromClientServer(server, &MockSyncCredentialsServer{})

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
		name         string
		data         []*workers.TextData
		expectedErr  error
		expectedBool bool
	}{
		{
			name:         "successful_sync",
			data:         []*workers.TextData{{ID: 1, UserID: 1, Text: "test", UpdatedAt: time.Now(), Action: "test", AccessToken: "test"}},
			expectedBool: true,
		},
		{
			name:         "data_is_nil",
			expectedBool: false,
		},
		{
			name:         "failed_sync",
			data:         []*workers.TextData{{ID: 1, UserID: 1, Text: "invalid", UpdatedAt: time.Now(), Action: "test", AccessToken: "test"}},
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandlers(logs)
			err = handler.SyncTextData(ctx, conn, tt.data)
			if (err != nil) == tt.expectedBool {
				t.Errorf("SyncCredentials() error = %v, wantErr %v", err, tt.expectedErr)

			}
		})
	}
}

func TestHandler_SyncBinaryData(t *testing.T) {
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
	v1_pd.RegisterSyncFromClientServer(server, &MockSyncCredentialsServer{})

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
		name         string
		data         []*workers.BinaryData
		expectedErr  error
		expectedBool bool
	}{
		{
			name:         "successful_sync",
			data:         []*workers.BinaryData{{ID: 1, UserID: 1, Binary: "test", UpdatedAt: time.Now(), Action: "test", AccessToken: "test"}},
			expectedBool: true,
		},
		{
			name:         "data_is_nil",
			expectedBool: false,
		},
		{
			name:         "failed_sync",
			data:         []*workers.BinaryData{{ID: 1, UserID: 1, Binary: "invalid", UpdatedAt: time.Now(), Action: "test", AccessToken: "test"}},
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandlers(logs)
			err = handler.SyncBinaryData(ctx, conn, tt.data)
			if (err != nil) == tt.expectedBool {
				t.Errorf("SyncCredentials() error = %v, wantErr %v", err, tt.expectedErr)

			}
		})
	}
}

func TestHandler_SyncCards(t *testing.T) {
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
	v1_pd.RegisterSyncFromClientServer(server, &MockSyncCredentialsServer{})

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
		name         string
		data         []*workers.Cards
		expectedErr  error
		expectedBool bool
	}{
		{
			name:         "successful_sync",
			data:         []*workers.Cards{{ID: 1, UserID: 1, Cards: "test", UpdatedAt: time.Now(), Action: "test", AccessToken: "test"}},
			expectedBool: true,
		},
		{
			name:         "data_is_nil",
			expectedBool: false,
		},
		{
			name:         "failed_sync",
			data:         []*workers.Cards{{ID: 1, UserID: 1, Cards: "invalid", UpdatedAt: time.Now(), Action: "test", AccessToken: "test"}},
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandlers(logs)
			err = handler.SyncCards(ctx, conn, tt.data)
			if (err != nil) == tt.expectedBool {
				t.Errorf("SyncCredentials() error = %v, wantErr %v", err, tt.expectedErr)

			}
		})
	}
}
