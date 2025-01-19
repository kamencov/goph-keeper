package binary_data

import (
	"context"
	"log/slog"
)


// storage - интерфейс сервиса.
//go:generate mockgen -source=service.go -destination=service_mock.go -package binary_data
type storage interface {
	GetUserIDByToken(ctx context.Context, accessToken string) (int, error)
	SaveBinaryDataBinary(ctx context.Context, uid int, data string) error
	DeletedBinary(ctx context.Context, userID int, data string) error
}


// Service - сервис для работы с данными.
type Service struct {
	log     *slog.Logger
	storage storage
}


// NewService - конструктор.
func NewService(log *slog.Logger, storage storage) *Service {
	return &Service{
		log:     log,
		storage: storage}
}


// SyncSaveBinary - сохраняет полученные данные.
func (s *Service) SyncSaveBinary(ctx context.Context, accessToken, data string) error {
	userID, err := s.storage.GetUserIDByToken(ctx, accessToken)
	if err != nil {
		s.log.Error("failed to get user_id")
		return err
	}
	err = s.storage.SaveBinaryDataBinary(ctx, userID, data)
	if err != nil {
		s.log.Error("failed to handlers data")
		return err
	}

	return nil
}


// SyncDelBinary - удаляет полученные данные.
func (s *Service) SyncDelBinary(ctx context.Context, accessToken, data string) error {
	userID, err := s.storage.GetUserIDByToken(ctx, accessToken)
	if err != nil {
		s.log.Error("failed to get user_id", "error", err)
		return err
	}

	if err := s.storage.DeletedBinary(ctx, userID, data); err != nil {
		s.log.Error("failed to deleted binary", "error", err)
		return err
	}

	return nil
}
