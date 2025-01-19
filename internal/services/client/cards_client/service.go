package cards_client

import (
	"context"
	"log/slog"
)

type storageClient interface {
	SaveCardsInDatabase(ctx context.Context, card string) error
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
