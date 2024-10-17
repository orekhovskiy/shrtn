package urlservice

func (s *Service) Ping() error {
	return s.urlRepository.Ping()
}
