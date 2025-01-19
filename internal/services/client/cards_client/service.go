package cards_client

import (
	"context"
	"log/slog"
)

// storageClient - представляет хранилище данных.
//go:generate mockgen -source=service.go -destination=service_mock.go -package=cards_client
type storageClient interface {
	SaveCardsInDatabase(ctx context.Context, userID int, card string) error
	SaveSync(ctx context.Context, tableName string, userID int, taskID int, action string) error
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
	GetIDTaskCards(ctx context.Context, tableName string, userID int, task string) (int, error)
}

// ServiceClient - представляет сервис для работы с cards.
type ServiceClient struct {
	log     *slog.Logger
	storage storageClient
}

// NewService - создает новый сервис клиента.
func NewService(log *slog.Logger, storage storageClient) *ServiceClient {
	return &ServiceClient{
		log:     log,
		storage: storage}
}

// SaveCards - сохраняет cards в базу данных.
func (s *ServiceClient) SaveCards(ctx context.Context, token, data string) error {
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		return err
	}

	if err = s.storage.SaveCardsInDatabase(ctx, userID, data); err != nil {
		s.log.Error("failed to save cards in database")
		return err
	}

	idTask, err := s.storage.GetIDTaskCards(ctx, "cards", userID, data)
	if err != nil {
		s.log.Error("failed to get id task")
		return err
	}

	if err = s.storage.SaveSync(ctx,
		"cards",
		userID,
		idTask,
		"save"); err != nil {
		s.log.Error("failed to save sync")
		return err
	}
	return nil
}
