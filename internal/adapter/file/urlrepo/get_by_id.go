package urlrepo

import (
	"fmt"
	"github.com/orekhovskiy/shrtn/internal/entity"
)

func (r *FileURLRepository) GetByID(id string) (*entity.URLRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exists := r.records[id]

	if !exists {
		return nil, fmt.Errorf("id not found: %s", id)
	}

	return &record, nil
}
