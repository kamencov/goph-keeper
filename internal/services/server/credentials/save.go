package credentials

import (
	"context"
	"errors"
)

var (
	ErrNotFoundUser = errors.New("user not found")
)

// SaveLoginAndPassword сохраняет логин и пароль от ресурса.
func (s *Service) SaveLoginAndPassword(ctx context.Context, userID int, resource, login, password string) error {

	err := s.storage.SaveLoginAndPasswordInCredentials(ctx, userID, resource, login, password)
	if err != nil {
		return err
	}

	return nil
}
