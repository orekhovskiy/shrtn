package urlrepo

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) GetByID(id string) (string, error) {
	var originalURL string
	err := r.db.QueryRow(getRecordById, id).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("id not found: %s", id)
		}
		return "", err
	}
	return originalURL, nil
}
