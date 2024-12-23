package cards_client

import "context"

func (s *ServiceClient) SaveCards(ctx context.Context, card string) error {
	return s.storage.SaveCardsInDatabase(ctx, card)
}
