package urlservice

func (s URLShortenerService) Save(url string, userID string) (string, error) {
	id := s.createShortURL(url)
	err := s.urlRepository.Save(id, url, userID)
	if err != nil {
		return "", err
	}

	return id, nil
}
