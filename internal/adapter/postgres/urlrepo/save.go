package urlrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	e "github.com/orekhovskiy/shrtn/internal/errors"
)

func (r *PostgresURLRepository) Save(id string, url string, userID string) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Try insert new record with conflict process
	query := `
		INSERT INTO url_records (uuid, short_url, original_url, user_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (original_url) DO NOTHING
		RETURNING short_url;
	`

	var shortURL string
	err = tx.QueryRow(query, uuid.New().String(), id, url, userID).Scan(&shortURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Conflict: record already exists, return short_url
			conflictQuery := `
				SELECT short_url FROM url_records
				WHERE original_url = $1;
			`
			var existingShortURL string
			queryErr := tx.QueryRow(conflictQuery, url).Scan(&existingShortURL)
			if queryErr != nil {
				return fmt.Errorf("unexpected error while retrieving conflicting URL: %w", queryErr)
			}

			return &e.URLConflictError{
				ShortURL:    existingShortURL,
				OriginalURL: url,
			}
		}

		// Internal error
		return fmt.Errorf("unexpected error while saving URL: %w", err)
	}

	// Finish transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
