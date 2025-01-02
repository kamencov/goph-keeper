package cards_client

import "context"

func (s *ServiceClient) SaveCards(ctx context.Context, token, data string) error {
	userID, err := s.storage.GetUserIDWithToken(ctx, token)
	if err != nil {
		return err
	}
	return s.storage.SaveCardsInDatabase(ctx, userID, data)
}
