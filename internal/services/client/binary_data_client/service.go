package binary_data_client

import (
	"context"
	"log/slog"
)

type storageClient interface {
	SaveBinaryDataInDatabase(ctx context.Context, userID int, data string) error
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
}

type ServiceClient struct {
	log     *slog.Logger
	storage storageClient
}

func NewService(log *slog.Logger, storage storageClient) *ServiceClient {
	return &ServiceClient{
		log:     log,
		storage: storage}
}
