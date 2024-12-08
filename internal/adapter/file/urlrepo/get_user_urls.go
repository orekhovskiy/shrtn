package urlrepo

import (
	"github.com/orekhovskiy/shrtn/internal/entity"
)

func (r *FileURLRepository) GetUserURLs(userID string) ([]entity.URLRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userURLs []entity.URLRecord
	for _, record := range r.records {
		if record.UserID == userID {
			userURLs = append(userURLs, record)
		}
	}

	return userURLs, nil
}
