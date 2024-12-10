package binary_data

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"goph-keeper/internal/middleware"
	pd "goph-keeper/internal/proto/v1"
	"log/slog"
)

// service - интерфейс сервисного слоя.
type service interface {
	SaveBinaryData(ctx context.Context, userID int, data string) error
}

// Handlers - структура ручки сохранения бинарных данных.
type Handlers struct {
	pd.UnimplementedPostBinaryDataServer
	service service
	log     *slog.Logger
}

// NewHandlers - конструктор ручки запроса сохранения в базу бинарные данные.
func NewHandlers(log *slog.Logger, service service) *Handlers {
	return &Handlers{
		log:     log,
		service: service,
	}
}

// PostBinaryData - обрабатывает запрос сохранения.
func (h *Handlers) PostBinaryData(ctx context.Context, in *pd.PostTextDataRequest) (*pd.Empty, error) {
	if in.Data == "" {
		h.log.Error("data is empty")
		return nil, status.Errorf(codes.InvalidArgument, "data is empty")
	}

	userID := ctx.Value(middleware.UserIDContextKey).(int)

	err := h.service.SaveBinaryData(ctx, userID, in.GetData())

	if err != nil {
		h.log.Error("failed to save login and password", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to save login and password")
	}

	return &pd.Empty{
		Message: "save complete",
	}, nil
}
