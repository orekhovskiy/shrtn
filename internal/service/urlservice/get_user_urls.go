package urlservice

import "github.com/orekhovskiy/shrtn/internal/entity"

func (s *URLShortenerService) GetUserURLs(userID string) ([]entity.URLRecord, error) {
	userURLs, err := s.urlRepository.GetUserURLs(userID)
	if err != nil {
		return nil, err
	}
	return userURLs, nil
}
