package credentials

import (
	"context"
	"errors"
)

var (
	ErrNotFoundUser = errors.New("user not found")
)

// SaveLoginAndPassword сохраняет логин и пароль от ресурса.
func (s *Service) SaveLoginAndPassword(ctx context.Context, resource, login, password string) error {

	loginID, err := s.storage.GetUserID(ctx, login)
	if err != nil {
		return ErrNotFoundUser
	}

	err = s.storage.SaveLoginAndPasswordInCredentials(ctx, resource, loginID, password)
	if err != nil {
		return err
	}

	return nil
}
