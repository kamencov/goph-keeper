package text_data

import (
	"context"
	"log/slog"
)

type storageTextData interface {
	SaveTextData(ctx context.Context, userID int, data string) error
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
