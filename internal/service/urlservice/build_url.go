package urlservice

import (
	"fmt"
	"net/url"
)

func (s *URLShortenerService) BuildURL(uri string) (string, error) {
	fullURL := s.options.BaseURL + "/" + uri
	_, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}
	return fullURL, nil
}
