package get_all_data

import (
	"context"
	"database/sql"
	"log/slog"
)

type storage interface {
	GetAll(ctx context.Context, tableName string) (*sql.Rows, error)
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
	return s.DB.GetAll(ctx, tableName)
}
