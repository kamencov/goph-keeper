package credentials

import (
	"context"
	"log/slog"
)

type credentials interface {
	SaveLoginAndPasswordInCredentials(ctx context.Context, resource string, loginID int, password string) error
	GetUserID(ctx context.Context, login string) (int, error)
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
