package credentials_client

import (
	"context"
	"errors"
	"log/slog"
)

// credentialsClient - интерфейс для работы с credentials.
//go:generate mockgen -source=service.go -destination=service_mock.go -package=credentials_client
type credentialsClient interface {
	SaveLoginAndPasswordInCredentials(ctx context.Context, userID int, resource, login, password string) error
	SaveSync(ctx context.Context, tableName string, userID int, taskID int, action string) error
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
	GetIDTaskCredentials(ctx context.Context, tableName string, userID int, task string) (int, error)
}

// ServiceClient - представляет сервис для работы с resources.
type ServiceClient struct {
	log     *slog.Logger
	storage credentialsClient
}

// NewService - создает новый экземпляр сервиса.
func NewService(log *slog.Logger, storage credentialsClient) *ServiceClient {
	return &ServiceClient{
		log:     log,
		storage: storage}
}

var (
	ErrNotFoundUser = errors.New("user not found")
)

// SaveLoginAndPassword сохраняет логин и пароль от ресурса.
func (s *ServiceClient) SaveLoginAndPassword(ctx context.Context, token, resource, login, password string) error {
	// получаем user_id
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		s.log.Error("failed to get user_id")
		return err
	}

	err = s.storage.SaveLoginAndPasswordInCredentials(ctx, userID, resource, login, password)
	if err != nil {
		s.log.Error("failed to handlers data")
		return err
	}

	idTask, err := s.storage.GetIDTaskCredentials(ctx, "credentials", userID, resource)
	if err != nil {
		s.log.Error("failed to get id task")
		return err
	}

	if err = s.storage.SaveSync(ctx, "credentials", userID, idTask, "save"); err != nil {
		s.log.Error("failed to save sync")
		return err
	}

	return nil
}
