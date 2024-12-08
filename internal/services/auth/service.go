package auth

import (
	"context"
	"log/slog"
	"time"
)

// storageAuth - интерфейс сервиса авторизации.
type storageAuth interface {
	CheckUser(ctx context.Context, login string) error
	CheckPassword(login string) (string, bool)
	SaveUser(ctx context.Context, login, hashPassword string) error
	SaveTableUserAndUpdateToken(login, accessToken string) error
	GetUserID(ctx context.Context, login string) (int, error)
}

// ServiceAuth - сервис авторизации.
type ServiceAuth struct {
	tokenSalt    []byte
	passwordSalt []byte

	accessTokenTTL time.Duration

	storage storageAuth

	log *slog.Logger
}

// NewServiceAuth - конструктор сервиса авторизации.
func NewServiceAuth(tokenSalt, passwordSalt []byte, accessTokenTTL time.Duration, log *slog.Logger, storage storageAuth) *ServiceAuth {
	return &ServiceAuth{
		tokenSalt:      tokenSalt,
		passwordSalt:   passwordSalt,
		accessTokenTTL: accessTokenTTL,
		log:            log,
		storage:        storage,
	}
}
