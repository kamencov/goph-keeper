package cards

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"goph-keeper/internal/middleware/auth"
	pd "goph-keeper/internal/proto/v1"
	"log/slog"
)

// service - интерфейс сервисного слоя.
type service interface {
	SaveCards(ctx context.Context, userID int, data string) error
}

// Handlers - структура ручки сохранения cards.
type Handlers struct {
	pd.UnimplementedPostCardsServer
	log     *slog.Logger
	service service
}

// NewHandlers - конструктор ручки запроса сохранения в базу данных cards.
func NewHandlers(log *slog.Logger, service service) *Handlers {
	return &Handlers{
		log:     log,
		service: service,
	}
}

// PostCards - обрабатывает запрос сохранения.
func (h *Handlers) PostCards(ctx context.Context, in *pd.PostTextDataRequest) (*pd.Empty, error) {

	if in.Data == "" {
		h.log.Error("data is empty")
		return nil, status.Errorf(codes.InvalidArgument, "data is empty")
	}

	userID := ctx.Value(auth.UserIDContextKey).(int)

	err := h.service.SaveCards(ctx, userID, in.GetData())

	if err != nil {
		h.log.Error("failed to handlers cards in base", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to handlers cards")
	}

	return &pd.Empty{
		Message: "handlers completed",
	}, nil
}
