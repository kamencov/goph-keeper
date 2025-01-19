package binary_data_client

import (
	"context"
	"log/slog"
)

type storageClient interface {
	SaveBinaryDataInDatabase(ctx context.Context, data string) error
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
