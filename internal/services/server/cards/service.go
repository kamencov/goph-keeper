package cards

import (
	"context"
	"log/slog"
)

// storageCards - интерфейс storage для сервиса Cards.
type storageCards interface {
	SaveCards(ctx context.Context, userID int, cards string) error
}

// ServiceCards - структура сервиса Cards.
type ServiceCards struct {
	log     *slog.Logger
	storage storageCards
}

// NewServiceCards - создаем сервис Cards, который заполняет структуру ServiceCards.
func NewServiceCards(log *slog.Logger, storage storageCards) *ServiceCards {
	return &ServiceCards{
		log:     log,
		storage: storage,
	}
}
