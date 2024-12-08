package urlservice

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"

	e "github.com/orekhovskiy/shrtn/internal/errors"
)

func (s *URLShortenerService) createShortURL(originalURL string) (string, error) {
	// Set initial random seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var shortURL string

	for {
		// Generate new salt
		salt := fmt.Sprintf("%d", r.Intn(1000000))

		// Make a hash with original url and generated salt
		hash := sha256.Sum256([]byte(originalURL + salt))
		shortURL = hex.EncodeToString(hash[:])[:7]

		// Check if collision occurred
		existingURL, err := s.GetByID(shortURL)
		if err != nil {
			if errors.Is(err, &e.NotFoundError{}) {
				// No collision occurred
				break
			}
			// Internal error
			return "", err
		}

		// If hashes collided, check original url
		if existingURL.OriginalURL == originalURL {
			return shortURL, nil
		}
	}

	// If no collision occurred, return generated short URL
	return shortURL, nil
}
