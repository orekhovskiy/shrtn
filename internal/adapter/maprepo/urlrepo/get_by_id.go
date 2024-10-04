package urlrepo

import (
	"fmt"
)

func (r *Repository) GetByID(id string) (string, error) {
	for _, record := range r.records {
		if record.ShortURL == id {
			return record.OriginalURL, nil
		}
	}
	return "", fmt.Errorf("id not found: %s", id)
}
