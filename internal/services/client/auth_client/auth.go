package auth_client

import (
	"context"
	"errors"
	"log/slog"
)

var (
	ErrPasswordNotCorrect = errors.New("password is not correct")
)

// storage - интерфейс хранилища данных.
//go:generate mockgen -source=auth.go -destination=auth_mock.go -package auth_client
type storage interface {
	SaveLoginAndToken(ctx context.Context, login, password, token string) error
	UpdateLoginAndToken(ctx context.Context, userID int, token string) error
	GetUserIDWithLogin(ctx context.Context, login string) (int, error)
	GetUserPassword(ctx context.Context, login string) (string, error)
	GetUserToken(ctx context.Context, login string) (string, error)
}

// Service - сервис авторизации.
type Service struct {
	log *slog.Logger
	db  storage
}

// NewService - конструктор.
func NewService(log *slog.Logger, db storage) *Service {
	return &Service{
		log: log,
		db:  db,
	}
}

// SaveTokenInBase - сохраняет токен в базе данных.
func (s *Service) SaveTokenInBase(ctx context.Context, login, password, token string) error {
	// получаем user_id с помощью login
	userID, err := s.db.GetUserIDWithLogin(ctx, login)
	if err != nil {
		err = s.db.SaveLoginAndToken(ctx, login, password, token)
		if err != nil {
			return err
		}
		return nil
	}
	err = s.db.UpdateLoginAndToken(ctx, userID, token)
	if err != nil {
		return err
	}

	return nil
}

// CheckUser - проверяет пользователя в базе данных.
func (s *Service) CheckUser(ctx context.Context, login, password string) (string, error) {
	pass, err := s.db.GetUserPassword(ctx, login)
	if err != nil {
		return "", ErrPasswordNotCorrect
	}

	if pass != password {
		return "", ErrPasswordNotCorrect
	}

	token, err := s.db.GetUserToken(ctx, login)
	if err != nil {
		return "", err
	}
	return token, nil
}
