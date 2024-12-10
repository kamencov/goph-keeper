package binary_data

import (
	"context"
	"log/slog"
)

type storage interface {
	SaveBinaryData(ctx context.Context, uid int, data string) error
}

type Service struct {
	log     *slog.Logger
	storage storage
}

func NewService(log *slog.Logger, storage storage) *Service {
	return &Service{
		log:     log,
		storage: storage}
}
