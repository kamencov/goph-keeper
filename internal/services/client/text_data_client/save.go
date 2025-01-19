package text_data_client

import "context"

func (s *ServiceClient) SaveTextData(ctx context.Context, data string) error {
	return s.storage.SaveTextDataInDatabase(ctx, data)
}
