package get_and_deleted_data

import (
	"context"
	"database/sql"
	"log/slog"
)


// storage - интерфейс для работы с базой данных.
//go:generate mockgen -source=service.go -destination=service_mock.go -package=get_and_deleted_data
type storage interface {
	SaveSync(ctx context.Context, tableName string, userID int, taskID int, action string) error
	GetAll(ctx context.Context, userID int, tableName string) (*sql.Rows, error)
	GetUserIDWithToken(ctx context.Context, token string) (int, error)
	Deleted(ctx context.Context, tableName string, id int) error
}

// GetAll - структура для работы с получением всех данных и удалением.
type GetAll struct {
	log *slog.Logger
	DB  storage
}

// NewService - конструктор.
func NewService(log *slog.Logger, db storage) *GetAll {
	return &GetAll{
		log: log,
		DB:  db,
	}
}

// GetAllData - возвращает все данные из базы данных.
func (s *GetAll) GetAllData(ctx context.Context, token, tableName string) (*sql.Rows, error) {
	userID, err := s.DB.GetUserIDWithToken(ctx, token)
	if err != nil {
		s.log.Error("failed to get user id with token", "error", err)
		return nil, err
	}

	return s.DB.GetAll(ctx, userID, tableName)
}

// DeletedData - удаляет данные из базы данных.
func (s *GetAll) DeletedData(ctx context.Context, token, tableName string, id int) error {
	userID, err := s.DB.GetUserIDWithToken(ctx, token)
	if err != nil {
		s.log.Error("failed to get user id with token", "error", err)
		return err
	}

	if err := s.DB.Deleted(ctx, tableName, id); err != nil {
		s.log.Error("failed to get user id with token")
		return err
	}

	if err = s.DB.SaveSync(ctx,
		tableName,
		userID,
		id,
		"deleted"); err != nil {
		s.log.Error("failed to save sync", "error", err)
		return err
	}
	return nil
}
