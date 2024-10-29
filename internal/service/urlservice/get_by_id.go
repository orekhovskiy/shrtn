package urlservice

func (s URLShortenerService) GetByID(id string) (string, error) {
	url, err := s.urlRepository.GetByID(id)

	if err != nil {
		return "", err
	}

	return url, nil
}
