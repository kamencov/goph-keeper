package credentials

func (s *Service) SaveLoginAndPassword(info, login, password string) error {

	loginID, err := s.storage.GetUserID(login)

	err := s.storage.SaveLoginAndPasswordInCredentials(info, login, password)
	if err != nil {
		return err
	}

	return nil
}
