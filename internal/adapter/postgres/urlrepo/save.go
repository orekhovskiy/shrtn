package urlrepo

import (
	"fmt"
	"github.com/google/uuid"
	e "github.com/orekhovskiy/shrtn/internal/errors"
)

func (r *PostgresURLRepository) Save(id string, url string, userID string) error {
	var existingShortURL string
	err := r.db.QueryRow(isRecordExists, url).Scan(&existingShortURL)
	if err == nil {
		return &e.URLConflictError{
			ShortURL:    existingShortURL,
			OriginalURL: url,
		}
	} else if err.Error() != "sql: no rows in result set" {
		return err
	}

	_, err = r.db.Exec(
		insertRecord,
		uuid.New().String(), id, url, userID)

	if err != nil {
		fmt.Println("error while inserting URL:", err)
		return err
	}

	return nil
}
