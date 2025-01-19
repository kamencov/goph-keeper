package credentials_client

import (
	"context"
	"errors"
)

var (
	ErrNotFoundUser = errors.New("user not found")
)

// SaveLoginAndPassword сохраняет логин и пароль от ресурса.
func (s *ServiceClient) SaveLoginAndPassword(ctx context.Context, token, resource, login, password string) error {
	// получаем user_id
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		s.log.Error("failed to get user_id")
		return err
	}

	err = s.storage.SaveLoginAndPasswordInCredentials(ctx, userID, resource, login, password)
	if err != nil {
		s.log.Error("failed to save data")
		return err
	}

	return nil
}
