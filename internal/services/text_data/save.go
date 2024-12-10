package text_data

import "context"

func (s *Service) SaveTextData(ctx context.Context, userID int, data string) error {
	return s.storage.SaveTextData(ctx, userID, data)
}
