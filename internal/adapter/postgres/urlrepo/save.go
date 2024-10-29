package urlrepo

import (
	"fmt"
	"github.com/google/uuid"
	e "github.com/orekhovskiy/shrtn/internal/errors"
)

func (r *PostgresURLRepository) Save(id string, url string) error {
	var existingShortURL string
	err := r.db.QueryRow(`SELECT short_url FROM url_records WHERE original_url = $1`, url).Scan(&existingShortURL)
	if err == nil {
		return &e.URLConflictError{
			ShortURL:    existingShortURL,
			OriginalURL: url,
		}
	} else if err.Error() != "sql: no rows in result set" {
		return err
	}

	_, err = r.db.Exec(
		`INSERT INTO url_records (uuid, short_url, original_url)
		 VALUES ($1, $2, $3)`,
		uuid.New().String(), id, url)

	if err != nil {
		fmt.Println("error while inserting URL:", err)
		return err
	}

	return nil
}
