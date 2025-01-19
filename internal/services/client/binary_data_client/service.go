package binary_data_client

import (
	"context"
	"log/slog"
)


// storageClient - представляет хранилище данных.
//go:generate mockgen -source=service.go -destination=service_mock.go -package=binary_data_client
type storageClient interface {
	SaveBinaryDataInDatabase(ctx context.Context, userID int, data string) error
	SaveSync(ctx context.Context, tableName string, userID int, taskID int, action string) error
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
	GetIDTaskBinary(ctx context.Context, tableName string, userID int, task string) (int, error)
}

// ServiceClient - представляет сервис для работы с бинарными данными.
type ServiceClient struct {
	log     *slog.Logger
	storage storageClient
}

// NewService - создает новый экземпляр ServiceClient.
func NewService(log *slog.Logger, storage storageClient) *ServiceClient {
	return &ServiceClient{
		log:     log,
		storage: storage}
}

// SaveBinaryData - сохраняет бинарные данные в базу данных.
func (s *ServiceClient) SaveBinaryData(ctx context.Context, token, data string) error {
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		s.log.Error("failed to get user id with token")
		return err
	}

	if err = s.storage.SaveBinaryDataInDatabase(ctx, userID, data); err != nil {
		s.log.Error("failed to save binary data in database")
		return err
	}

	idTask, err := s.storage.GetIDTaskBinary(ctx, "binary_data", userID, data)
	if err != nil {
		s.log.Error("failed to get id task")
		return err
	}

	if err = s.storage.SaveSync(ctx, "binary_data", userID, idTask, "save"); err != nil {
		s.log.Error("failed to save sync")
		return err
	}

	return nil
}
