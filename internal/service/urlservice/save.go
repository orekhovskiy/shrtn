package urlservice

import (
	"crypto/sha256"
	"encoding/hex"
)

func (s URLShortenerService) Save(url string) (string, error) {
	hash := sha256.Sum256([]byte(url))
	id := hex.EncodeToString(hash[:])[:7]
	err := s.urlRepository.Save(id, url)
	if err != nil {
		return "", err
	}

	return id, nil
}
