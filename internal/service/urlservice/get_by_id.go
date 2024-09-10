package urlservice

func (s Service) GetById(id string) (string, error) {
	url, err := s.urlRepository.GetById(id)

	if err != nil {
		return "", err
	}

	return url, nil
}
