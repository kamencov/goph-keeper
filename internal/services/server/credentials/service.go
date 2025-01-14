package credentials

import (
	"context"
	"errors"
	"log/slog"
)

type credentials interface {
	GetUserIDByToken(ctx context.Context, accessToken string) (int, error)
	ServerSaveLoginAndPasswordInCredentials(ctx context.Context, userID int, resource, login, password string) error
	DeletedCredentials(ctx context.Context, userID int, resource string) error
}

type Service struct {
	log     *slog.Logger
	storage credentials
}

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

func (s *Service) SyncDelCredentials(ctx context.Context, accessToken, resource string) error {
	userID, err := s.storage.GetUserIDByToken(ctx, accessToken)
	if err != nil {
		s.log.Error("failed to get user_id", "error", err)
		return err
	}

	if err := s.storage.DeletedCredentials(ctx, userID, resource); err != nil {
		s.log.Error("failed to deleted credentials", "error", err)
		return err
	}

	return nil
}
