package urlservice

import (
	"crypto/sha256"
	"encoding/hex"
)

func (s *URLShortenerService) createShortURL(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL))
	shortURL := hex.EncodeToString(hash[:])[:7]
	return shortURL
}
