package binary_data_client

import "context"

func (s *ServiceClient) SaveBinaryData(ctx context.Context, data string) error {

	return s.storage.SaveBinaryDataInDatabase(ctx, data)
}
