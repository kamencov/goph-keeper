package health

import (
	"golang.org/x/net/context"
	pd "goph-keeper/internal/proto/v1"
	"log/slog"
)

type Handler struct {
	pd.UnimplementedHealthServer
	log *slog.Logger
}

func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) Health(ctx context.Context, in *pd.Empty) (*pd.Empty, error) {
	return &pd.Empty{
		Message: "SERVING",
	}, nil
}
