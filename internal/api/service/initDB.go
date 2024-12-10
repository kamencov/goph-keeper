package service

import (
	"context"
	"fmt"
	"goph-keeper/internal/storage/memory"
	"goph-keeper/internal/storage/postgresql"
	"log/slog"
)

type repo interface {
	CheckUser(ctx context.Context, login string) error
	CheckPassword(login string) (string, bool)
	SaveUser(ctx context.Context, login, hashPassword string) error
	SaveTableUserAndUpdateToken(login, accessToken string) error
	SaveLoginAndPasswordInCredentials(ctx context.Context, resource string, loginID int, password string) error
	GetUserID(ctx context.Context, login string) (int, error)
	Close() error
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
