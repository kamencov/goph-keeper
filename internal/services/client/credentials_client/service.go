package credentials_client

import (
	"context"
	"log/slog"
)

type credentialsClient interface {
	SaveLoginAndPasswordInCredentials(ctx context.Context, userID int, resource, login, password string) error
	SaveSync(ctx context.Context, tableName string, userID int, taskID int, action string) error
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
	GetIDTaskCredentials(ctx context.Context, tableName string, userID int, task string) (int, error)
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
