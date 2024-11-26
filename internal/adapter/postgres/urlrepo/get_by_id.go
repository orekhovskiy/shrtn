package urlrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/orekhovskiy/shrtn/internal/entity"
)

func (r *PostgresURLRepository) GetByID(id string) (*entity.URLRecord, error) {
	record := &entity.URLRecord{}
	err := r.db.QueryRow(getRecordByID, id).Scan(
		&record.UUID,
		&record.ShortURL,
		&record.OriginalURL,
		&record.IsDeleted,
		&record.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("id not found: %s", id)
		}
		return nil, err
	}
	return record, nil
}
