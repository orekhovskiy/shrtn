package urlrepo

import (
	"fmt"
	"slices"
)

func (r Repository) GetByID(id string) (string, error) {
	index := slices.IndexFunc(r.records, func(record URLRecord) bool {
		return record.ShortURL == id
	})

	if index == -1 {
		return "", fmt.Errorf("id not found: %s", id)
	}
	return r.records[index].OriginalURL, nil
}
