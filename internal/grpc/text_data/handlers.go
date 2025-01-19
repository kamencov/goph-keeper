package text_data

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
	SaveTextData(ctx context.Context, userID int, data string) error
}

// Handlers - структура ручки сохранения текста.
type Handlers struct {
	pd.UnimplementedPostTextDataServer
	log     *slog.Logger
	service service
}

// NewHandlers - конструктор ручки запроса сохранения в базу данных текста.
func NewHandlers(log *slog.Logger, service service) *Handlers {
	return &Handlers{
		log:     log,
		service: service,
	}
}

// PostTextData - обрабатывает запрос сохранения.
func (h *Handlers) PostTextData(ctx context.Context, in *pd.PostTextDataRequest) (*pd.Empty, error) {
	if in.Data == "" {
		h.log.Error("data is empty")
		return nil, status.Errorf(codes.InvalidArgument, "data is empty")
	}

	userID := ctx.Value(middleware.UserIDContextKey).(int)

	err := h.service.SaveTextData(ctx, userID, in.GetData())

	if err != nil {
		h.log.Error("failed to save login and password", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to save login and password")
	}

	return &pd.Empty{
		Message: "save complete",
	}, nil
}
