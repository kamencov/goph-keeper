package credentials_client

import (
	"context"
	"goph-keeper/internal/api/client/cli"
)

func (s *ServiceClient) GetResourceLoginAndPassword(ctx context.Context, token string) (*[]cli.Resource, error) {
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		s.log.Error("failed to get user id", "error", err)
		return nil, err
	}

	result, err := s.storage.GetLoginAndPasswordInCredentials(ctx, userID)
	if err != nil {
		s.log.Error("failed to get credentials data", "error", err)
		return nil, err
	}

	return result, nil
}
