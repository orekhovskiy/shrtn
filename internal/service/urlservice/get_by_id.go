package urlservice

func (s Service) GetByID(id string) (string, error) {
	url, err := s.urlRepository.GetByID(id)

	if err != nil {
		return "", err
	}

	return url, nil
}
