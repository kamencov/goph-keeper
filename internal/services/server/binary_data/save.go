package binary_data

import "context"

func (s *Service) SaveBinaryData(ctx context.Context, userID int, data string) error {

	return s.storage.SaveBinaryData(ctx, userID, data)
}
