package urlrepo

import (
	"fmt"
)

func (r Repository) GetById(id string) (string, error) {
	url, exists := r.urlMapping[id]
	if !exists {
		return "", fmt.Errorf("id not found: %s", id)
	}
	return url, nil
}
