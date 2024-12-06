package urlservice

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	e "github.com/orekhovskiy/shrtn/internal/errors"
)

func (s *URLShortenerService) createShortURL(originalURL string) (string, error) {
	// Generate the hash for the original URL and shorten it to 7 characters
	hash := sha256.Sum256([]byte(originalURL))
	shortURL := hex.EncodeToString(hash[:])[:7]

	// Check if the short URL already exists
	existingURL, err := s.GetByID(shortURL)
	if err != nil {
		if errors.Is(err, &e.NotFoundError{}) {
			// No collision, it's a new URL
			return shortURL, nil
		}
		// If the error is something else, return it
		return "", err
	}

	// If we got here, it means the short URL already maps to an existing original URL
	// Therefore, it's a collision
	if existingURL.OriginalURL != originalURL {
		return "", fmt.Errorf(
			"collided URL: the short URL '%s' already maps to a different original URL '%s'",
			shortURL,
			existingURL.OriginalURL,
		)
	}

	// If no collision, return the existing short URL
	return shortURL, nil
}
