package text_data_client

import (
	"context"
	"log/slog"
)
// storageTextDataClient - интерфейс для работы с текстовыми данными.
//go:generate mockgen -source=service.go -destination=service_mock.go -package=text_data_client
type storageTextDataClient interface {
	SaveTextDataInDatabase(ctx context.Context, userID int, data string) error
	SaveSync(ctx context.Context, tableName string, userID int, taskID int, action string) error
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
	GetIDTaskText(ctx context.Context, tableName string, userID int, task string) (int, error)
}

// ServiceClient - сервис для работы с текстовыми данными.
type ServiceClient struct {
	log     *slog.Logger
	storage storageTextDataClient
}

// NewService - конструктор сервиса для работы с текстовыми данными.
func NewService(log *slog.Logger, storage storageTextDataClient) *ServiceClient {
	return &ServiceClient{
		log:     log,
		storage: storage}
}

// SaveTextData - сохраняет текстовые данные.
func (s *ServiceClient) SaveTextData(ctx context.Context, token, data string) error {
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		s.log.Error("failed to get user id with token", "error", err)
		return err
	}

	if err = s.storage.SaveTextDataInDatabase(ctx, userID, data); err != nil {
		s.log.Error("failed to save text data in database", "error", err)
		return err
	}

	idTask, err := s.storage.GetIDTaskText(ctx, "text_data", userID, data)
	if err != nil {
		s.log.Error("failed to get id task", "error", err)
		return err
	}

	if err = s.storage.SaveSync(ctx,
		"text_data",
		userID,
		idTask,
		"save"); err != nil {
		s.log.Error("failed to save sync", "error", err)
		return err
	}

	return nil
}
