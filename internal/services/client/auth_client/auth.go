package auth_client

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

type storage interface {
	SaveLoginAndToken(ctx context.Context, login, token string) error
	UpdateLoginAndToken(ctx context.Context, userID int, token string) error
	GetUserIDWithLogin(ctx context.Context, login string) (int, error)
}

type Service struct {
	log *slog.Logger
	db  storage
}

func NewService(log *slog.Logger, db storage) *Service {
	return &Service{
		log: log,
		db:  db,
	}
}

func (s *Service) SaveTokenInBase(ctx context.Context, login, token string) error {
	// получаем user_id с помощью login
	userID, err := s.db.GetUserIDWithLogin(ctx, login)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			err = s.db.SaveLoginAndToken(ctx, login, token)
			if err != nil {
				return err
			}
			return nil
		}
		s.log.Error("failed to check user id", "error:", err)
		return err
	}
	err = s.db.UpdateLoginAndToken(ctx, userID, token)
	if err != nil {
		return err
	}

	return nil
}
