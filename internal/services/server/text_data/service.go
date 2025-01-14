package text_data

import (
	"context"
	"log/slog"
)

type storageTextData interface {
	GetUserIDByToken(ctx context.Context, accessToken string) (int, error)
	SaveTextDataPstgres(ctx context.Context, userID int, data string) error
	DeletedText(ctx context.Context, userID int, data string) error
}

type Service struct {
	log     *slog.Logger
	storage storageTextData
}

func NewService(log *slog.Logger, storage storageTextData) *Service {
	return &Service{
		log:     log,
		storage: storage}
}

func (s *Service) SaveTextData(ctx context.Context, userID int, data string) error {
	return s.storage.SaveTextDataPstgres(ctx, userID, data)
}

func (s *Service) SyncSaveText(ctx context.Context, accessToken, data string) error {
	userID, err := s.storage.GetUserIDByToken(ctx, accessToken)
	if err != nil {
		s.log.Error("failed to get user_id", "error", err)
		return err
	}
	err = s.storage.SaveTextDataPstgres(ctx, userID, data)
	if err != nil {
		s.log.Error("failed to handlers data", "error", err)
		return err
	}

	return nil
}

func (s *Service) SyncDelText(ctx context.Context, accessToken, data string) error {
	userID, err := s.storage.GetUserIDByToken(ctx, accessToken)
	if err != nil {
		s.log.Error("failed to get user_id", "error", err)
		return err
	}

	if err := s.storage.DeletedText(ctx, userID, data); err != nil {
		s.log.Error("failed to deleted text", "error", err)
		return err
	}

	return nil
}
