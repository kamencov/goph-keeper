package service

import (
	"context"
	"fmt"
	"goph-keeper/internal/storage/memory"
	"goph-keeper/internal/storage/postgresql"
	"log/slog"
)

type repo interface {
	Close() error
	CheckUser(ctx context.Context, login string) error
	CheckPassword(login string) (string, bool)
	SaveUser(ctx context.Context, login, hashPassword string) error
	SaveTableUserAndUpdateToken(login, accessToken string) error
	SaveLoginAndPasswordInCredentials(ctx context.Context, userID int, resource string, login string, password string) error
	SaveTextData(ctx context.Context, userID int, data string) error
	SaveBinaryData(ctx context.Context, uid int, data string) error
	SaveCards(ctx context.Context, userID int, cards string) error
	GetUserIDByToken(ctx context.Context, token string) (int, error)
	GetUserIDByLogin(ctx context.Context, login string) (int, error)
}

func initDB(log *slog.Logger, f *Flags) (repo, error) {

	switch f.Repo {
	case "1":
		return memory.NewMemory(log), nil
	case "2":
		return postgresql.NewPostgresql(log)
	default:
		return nil, fmt.Errorf("unknown repository type: %s", f.Repo)
	}
}
