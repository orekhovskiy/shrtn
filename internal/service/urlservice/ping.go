package urlservice

func (s *URLShortenerService) Ping() error {
	return s.urlRepository.Ping()
}
