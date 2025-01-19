package health

import (
	"golang.org/x/net/context"
	pd "goph-keeper/internal/proto/v1"
	"log/slog"
)

// Handler - обработчик запросов.
type Handler struct {
	pd.UnimplementedHealthServer
	log *slog.Logger
}

// NewHandler - конструктор обработчика.
func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

// Health - проверяет работоспособность сервиса.
func (h *Handler) Health(ctx context.Context, in *pd.Empty) (*pd.Empty, error) {
	return &pd.Empty{
		Message: "SERVING",
	}, nil
}
