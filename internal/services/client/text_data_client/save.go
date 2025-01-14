package text_data_client

import "context"

func (s *ServiceClient) SaveTextData(ctx context.Context, token, data string) error {
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		s.log.Error("failed to get user id with token", "error", err)
		return err
	}

	if err = s.storage.SaveTextDataInDatabase(ctx, userID, data); err != nil {
		s.log.Error("failed to save text data in database", "error", err)
		return err
	}

	idTask, err := s.storage.GetIDTaskText(ctx, "text_data", userID, data)
	if err != nil {
		s.log.Error("failed to get id task", "error", err)
		return err
	}

	if err = s.storage.SaveSync(ctx,
		"text_data",
		userID,
		idTask,
		"save"); err != nil {
		s.log.Error("failed to save sync", "error", err)
		return err
	}

	return nil
}
