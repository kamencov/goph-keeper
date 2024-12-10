package text_data

import "log/slog"

type storageTextData interface {
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
