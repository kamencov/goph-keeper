package health

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	v1_pd "goph-keeper/internal/proto/v1"
)

// Handler - интерфейс на сервисный слой.
type Handler struct {
	log *slog.Logger
}

// NewHandlers - конструктор обработчика.
func NewHandlers(log *slog.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

// Health - проверяет работоспособность сервиса.
func (h *Handler) Health(ctx context.Context, conn *grpc.ClientConn) error {
	healthClient := v1_pd.NewHealthClient(conn)

	_, err := healthClient.Health(ctx, &v1_pd.Empty{})

	if err != nil {
		h.log.Error("health.repositories.app: failed to register user")
		return err
	}
	return nil
}
