package cards

import (
	"context"
	"log/slog"
)

// storageCards - интерфейс storage для сервиса Cards.
//go:generate mockgen -source=service.go -destination=service_mock.go -package cards
type storageCards interface {
	GetUserIDByToken(ctx context.Context, accessToken string) (int, error)
	SaveCards(ctx context.Context, userID int, cards string) error
	DeletedCards(ctx context.Context, userID int, data string) error
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

// SyncSaveCards - сохраняет полученные данные.
func (s *ServiceCards) SyncSaveCards(ctx context.Context, accessToken, cards string) error {
	userID, err := s.storage.GetUserIDByToken(ctx, accessToken)
	if err != nil {
		s.log.Error("failed to get user_id")
		return err
	}
	err = s.storage.SaveCards(ctx, userID, cards)
	if err != nil {
		s.log.Error("failed to handlers data")
		return err
	}

	return nil
}

// SyncDelBinary - удаляет полученные данные.
func (s *ServiceCards) SyncDelBinary(ctx context.Context, accessToken, data string) error {
	userID, err := s.storage.GetUserIDByToken(ctx, accessToken)
	if err != nil {
		s.log.Error("failed to get user_id", "error", err)
		return err
	}

	if err := s.storage.DeletedCards(ctx, userID, data); err != nil {
		s.log.Error("failed to deleted cards", "error", err)
		return err
	}

	return nil
}
