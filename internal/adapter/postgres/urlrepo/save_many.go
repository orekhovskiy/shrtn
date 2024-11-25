package urlrepo

import (
	"github.com/google/uuid"
	"github.com/orekhovskiy/shrtn/internal/entity"
	"strconv"
	"strings"
)

func (r *PostgresURLRepository) SaveMany(records []entity.URLRecord, userID string) ([]entity.URLRecord, error) {
	if len(records) == 0 {
		return nil, nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	insertedRecords := make([]entity.URLRecord, 0, len(records))

	var values []string
	var args []interface{}
	for i, record := range records {
		if record.UUID == "" {
			record.UUID = uuid.New().String()
		}
		record.UserID = userID

		values = append(values,
			"($"+strconv.Itoa(i*4+1)+", $"+strconv.Itoa(i*4+2)+", $"+strconv.Itoa(i*4+3)+", $"+strconv.Itoa(i*4+4)+")")
		args = append(args, record.UUID, record.ShortURL, record.OriginalURL, record.UserID)
	}

	sqlQuery := `
        INSERT INTO url_records (uuid, short_url, original_url, user_id)
        VALUES ` + strings.Join(values, ", ") + `
        ON CONFLICT (short_url) DO NOTHING
        RETURNING uuid, short_url, original_url
    `

	rows, err := tx.Query(sqlQuery, args...)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	for rows.Next() {
		var insertedRecord entity.URLRecord
		if err := rows.Scan(&insertedRecord.UUID, &insertedRecord.ShortURL, &insertedRecord.OriginalURL, nil); err != nil {
			err := tx.Rollback()
			if err != nil {
				return nil, err
			}
			return nil, err
		}
		insertedRecords = append(insertedRecords, insertedRecord)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return insertedRecords, nil
}
