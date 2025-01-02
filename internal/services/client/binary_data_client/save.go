package binary_data_client

import "context"

func (s *ServiceClient) SaveBinaryData(ctx context.Context, token, data string) error {
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		return err
	}

	return s.storage.SaveBinaryDataInDatabase(ctx, userID, data)
}
