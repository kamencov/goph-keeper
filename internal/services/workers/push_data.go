package workers

type credentials struct {
}

func (s *Service) PushData() {
	rows, err := s.storage.GetAllNewCredentials()
	if err != nil {
		s.log.Error("failed to get data from database", "error", err)
		return
	}
	s.log.Info("got data from database", "rows", rows)

	for rows.Next() {

	}
}
