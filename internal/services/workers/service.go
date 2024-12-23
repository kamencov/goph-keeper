package workers

import (
	"database/sql"
	"log/slog"
)

type storage interface {
	GetAllNewCredentials() (*sql.Rows, error)
	GetAllNewTextData() (*sql.Rows, error)
	GetAllNewBinaryData() (*sql.Rows, error)
	GetAllNewCards() (*sql.Rows, error)
}
type Service struct {
	log     *slog.Logger
	storage storage
}

func NewService(log *slog.Logger, storage storage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}
