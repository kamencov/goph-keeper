package get_and_deleted_data

import (
	"context"
	"database/sql"
	"log/slog"
)

type storage interface {
	GetAll(ctx context.Context, userID int, tableName string) (*sql.Rows, error)
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
	Deleted(ctx context.Context, tableName string, id int) error
}

type GetAll struct {
	log *slog.Logger
	DB  storage
}

func NewService(log *slog.Logger, db storage) *GetAll {
	return &GetAll{
		log: log,
		DB:  db,
	}
}

func (s *GetAll) GetAllData(ctx context.Context, token, tableName string) (*sql.Rows, error) {
	userID, err := s.DB.GetUserIDWithToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return s.DB.GetAll(ctx, userID, tableName)
}

func (s *GetAll) DeletedData(ctx context.Context, tableName string, id int) error {
	return s.DB.Deleted(ctx, tableName, id)
}
