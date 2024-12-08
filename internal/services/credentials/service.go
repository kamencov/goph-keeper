package credentials

import "log/slog"

type credentials interface {
	SaveLoginAndPasswordInCredentials(info, login, password string) error
}

type Service struct {
	log     *slog.Logger
	storage credentials
}

func NewService(log *slog.Logger, storage credentials) *Service {
	return &Service{
		log:     log,
		storage: storage}
}
