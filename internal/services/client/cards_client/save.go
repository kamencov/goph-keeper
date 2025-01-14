package cards_client

import "context"

func (s *ServiceClient) SaveCards(ctx context.Context, token, data string) error {
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		return err
	}

	if err = s.storage.SaveCardsInDatabase(ctx, userID, data); err != nil {
		s.log.Error("failed to save cards in database")
		return err
	}

	idTask, err := s.storage.GetIDTaskCards(ctx, "cards", userID, data)
	if err != nil {
		s.log.Error("failed to get id task")
		return err
	}

	if err = s.storage.SaveSync(ctx,
		"cards",
		userID,
		idTask,
		"save"); err != nil {
		s.log.Error("failed to save sync")
		return err
	}
	return nil
}
