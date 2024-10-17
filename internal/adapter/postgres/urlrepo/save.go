package urlrepo

import (
	"github.com/google/uuid"
)

func (r *Repository) Save(id string, url string) error {
	var exists bool
	err := r.db.QueryRow(isRecordExists, id).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = r.db.Exec(
		insertRecord,
		uuid.New().String(), id, url,
	)
	return err
}
