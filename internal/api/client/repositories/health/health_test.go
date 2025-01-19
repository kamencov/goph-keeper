package health

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	v1_pd "goph-keeper/internal/proto/v1"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"
)

// MockHelthServer - моковая реализация сервера Health.
type MockHealthServer struct {
	v1_pd.UnimplementedHealthServer
}

// Health - имитация проверки работоспособности сервиса
func (*MockHealthServer) Health(ctx context.Context, in *v1_pd.Empty) (*v1_pd.Empty, error) {
	return &v1_pd.Empty{}, nil
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

func TestHandler_Health(t *testing.T) {
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
	v1_pd.RegisterHealthServer(server, &MockHealthServer{})

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
		in           v1_pd.Empty
		expectedBool bool
	}{
		{
			name: "Health",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewHandlers(logs)
			err = h.Health(ctx, conn)

			if err != nil {
				t.Errorf("Health() error = %v, wantErr %v", err, test.expectedBool)
			}
		})
	}
}
