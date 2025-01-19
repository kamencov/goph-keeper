package credentials

import (
	"context"
	"errors"
	"log/slog"
)

// credentials - интерфейс сервиса для работы с данными.
//
//go:generate mockgen -source=service.go -destination=service_mock.go -package credentials
type credentials interface {
	GetUserIDByToken(ctx context.Context, accessToken string) (int, error)
	ServerSaveLoginAndPasswordInCredentials(ctx context.Context, userID int, resource, login, password string) error
	DeletedCredentials(ctx context.Context, userID int, resource string) error
}

// Service - сервис для работы с данными.
type Service struct {
	log     *slog.Logger
	storage credentials
}

// NewService - создает новый экземпляр сервиса.
func NewService(log *slog.Logger, storage credentials) *Service {
	return &Service{
		log:     log,
		storage: storage}
}

var (
	ErrNotFoundUser = errors.New("user not found")
)

// SyncSaveCredentials сохраняет логин и пароль от ресурса.
func (s *Service) SyncSaveCredentials(ctx context.Context, accessToken, resource, login, password string) error {
	userID, err := s.storage.GetUserIDByToken(ctx, accessToken)
	if err != nil {
		s.log.Error("failed to get user_id", "error", err)
		return ErrNotFoundUser
	}
	err = s.storage.ServerSaveLoginAndPasswordInCredentials(ctx, userID, resource, login, password)
	if err != nil {
		s.log.Error("failed to handlers data", "error", err)
		return err
	}

	return nil
}

// SyncDelCredentials удаляет логин и пароль от ресурса.
func (s *Service) SyncDelCredentials(ctx context.Context, accessToken, resource string) error {
	userID, err := s.storage.GetUserIDByToken(ctx, accessToken)
	if err != nil {
		s.log.Error("failed to get user_id", "error", err)
		return ErrNotFoundUser
	}

	if err := s.storage.DeletedCredentials(ctx, userID, resource); err != nil {
		s.log.Error("failed to deleted credentials", "error", err)
		return err
	}

	return nil
}
