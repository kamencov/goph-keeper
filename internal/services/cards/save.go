package cards

import "context"

// SaveCards - отрабатывает полученные данные в слой storage.
func (s *ServiceCards) SaveCards(ctx context.Context, userID int, cards string) error {
	return s.storage.SaveCards(ctx, userID, cards)
}
