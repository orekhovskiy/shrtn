package urlrepo

import (
	"database/sql"
	"errors"

	"github.com/orekhovskiy/shrtn/internal/entity"
	e "github.com/orekhovskiy/shrtn/internal/errors"
)

func (r *PostgresURLRepository) GetByID(id string) (*entity.URLRecord, error) {
	record := &entity.URLRecord{}
	getRecordByID := `
        SELECT uuid, short_url, original_url, is_deleted, user_id
        FROM url_records
        WHERE short_url = $1
        LIMIT 1;
    `
	err := r.db.QueryRow(getRecordByID, id).Scan(
		&record.UUID,
		&record.ShortURL,
		&record.OriginalURL,
		&record.IsDeleted,
		&record.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &e.NotFoundError{ID: id}
		}
		return nil, err
	}

	return record, nil
}
