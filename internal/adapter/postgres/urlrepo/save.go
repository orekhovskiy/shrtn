package urlrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	e "github.com/orekhovskiy/shrtn/internal/errors"
)

const pgUniqueConstraintViolationCode = "23505"

func (r *PostgresURLRepository) Save(id string, url string, userID string) error {
	// Try insert new record with conflict process
	query := `
		INSERT INTO url_records (uuid, short_url, original_url, user_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (original_url) DO UPDATE
		SET short_url = EXCLUDED.short_url
		RETURNING short_url;
	`

	var existingShortURL string
	err := r.db.QueryRow(query, uuid.New().String(), id, url, userID).Scan(&existingShortURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No conflict, record successfully added
			return nil
		}

		// Ensure unique constraint violated
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueConstraintViolationCode {
			return &e.URLConflictError{
				ShortURL:    existingShortURL,
				OriginalURL: url,
			}
		}

		// Internal error
		return fmt.Errorf("unexpected error while saving URL: %w", err)
	}

	return nil
}
