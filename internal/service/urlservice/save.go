package urlservice

func (s URLShortenerService) Save(url string) (string, error) {
	id := s.createShortURL(url)
	err := s.urlRepository.Save(id, url)
	if err != nil {
		return "", err
	}

	return id, nil
}
