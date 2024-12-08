package urlrepo

import (
	"github.com/orekhovskiy/shrtn/internal/entity"
	e "github.com/orekhovskiy/shrtn/internal/errors"
)

func (r *FileURLRepository) GetByID(id string) (*entity.URLRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exists := r.records[id]

	if !exists {
		return nil, &e.NotFoundError{ID: id}
	}

	return &record, nil
}
