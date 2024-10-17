package urlrepo

import (
	"fmt"
)

func (r *Repository) GetByID(id string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exists := r.records[id]
	if !exists {
		return "", fmt.Errorf("id not found: %s", id)
	}
	return record.OriginalURL, nil
}
