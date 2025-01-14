package binary_data_client

import "context"

func (s *ServiceClient) SaveBinaryData(ctx context.Context, token, data string) error {
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		s.log.Error("failed to get user id with token")
		return err
	}

	if err = s.storage.SaveBinaryDataInDatabase(ctx, userID, data); err != nil {
		s.log.Error("failed to save binary data in database")
		return err
	}

	idTask, err := s.storage.GetIDTaskBinary(ctx, "binary_data", userID, data)
	if err != nil {
		s.log.Error("failed to get id task")
		return err
	}

	if err = s.storage.SaveSync(ctx, "binary_data", userID, idTask, "save"); err != nil {
		s.log.Error("failed to save sync")
		return err
	}

	return nil
}
