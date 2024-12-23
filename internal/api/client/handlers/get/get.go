package get

import (
	"context"
	"goph-keeper/internal/api/client/cli"
	"log/slog"
)

type serviceResource interface {
	GetResourceLoginAndPassword(ctx context.Context, token string) (*[]cli.Resource, error)
}

type Handler struct {
	log     *slog.Logger
	service serviceResource
}

func NewHandler(log *slog.Logger, service serviceResource) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) GetResource(ctx context.Context, token string) (*[]cli.Resource, error) {
	result, err := h.service.GetResourceLoginAndPassword(ctx, token)
	if err != nil {
		h.log.Error("failed to get resource data", "error", err)
		return nil, err
	}

	return result, nil
}
