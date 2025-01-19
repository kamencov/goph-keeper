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
// @Tags GET
// @Summary Проверяет работоспособность сервиса.
// @Description Проверяет работоспособность сервиса.
// @Accept json
// @Produce json
// @Success 200 {object} v1_pd.Empty
// @Router /goph_keeper_v1.Health/Health [get]
func (h *Handler) Health(ctx context.Context, in *pd.Empty) (*pd.Empty, error) {
	return &pd.Empty{
		Message: "SERVING",
	}, nil
}
