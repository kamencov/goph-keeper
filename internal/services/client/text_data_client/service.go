package text_data_client

import (
	"context"
	"log/slog"
)

type storageTextDataClient interface {
	SaveTextDataInDatabase(ctx context.Context, userID int, data string) error
	SaveSync(ctx context.Context, tableName string, userID int, taskID int, action string) error
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
	GetIDTaskText(ctx context.Context, tableName string, userID int, task string) (int, error)
}

type ServiceClient struct {
	log     *slog.Logger
	storage storageTextDataClient
}

func NewService(log *slog.Logger, storage storageTextDataClient) *ServiceClient {
	return &ServiceClient{
		log:     log,
		storage: storage}
}
