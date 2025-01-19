package text_data_client

import (
	"context"
	"log/slog"
)

type storageTextDataClient interface {
	SaveTextDataInDatabase(ctx context.Context, data string) error
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
