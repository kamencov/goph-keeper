package credentials_client

import (
	"context"
	"goph-keeper/internal/api/client/cli"
	"log/slog"
)

type credentialsClient interface {
	SaveLoginAndPasswordInCredentials(ctx context.Context, userID int, resource, login, password string) error
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
	GetLoginAndPasswordInCredentials(ctx context.Context, userID int) (*[]cli.Resource, error)
}

type ServiceClient struct {
	log     *slog.Logger
	storage credentialsClient
}

func NewService(log *slog.Logger, storage credentialsClient) *ServiceClient {
	return &ServiceClient{
		log:     log,
		storage: storage}
}
