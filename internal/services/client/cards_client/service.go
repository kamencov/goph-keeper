package cards_client

import (
	"context"
	"log/slog"
)

type storageClient interface {
	SaveCardsInDatabase(ctx context.Context, userID int, card string) error
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
