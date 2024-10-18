package urlrepo

import (
	"github.com/google/uuid"
	"github.com/orekhovskiy/shrtn/internal/entity"
	"strconv"
	"strings"
)

func (r *Repository) SaveMany(records []entity.URLRecord) ([]entity.URLRecord, error) {
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

		values = append(values,
			"($"+strconv.Itoa(i*3+1)+", $"+strconv.Itoa(i*3+2)+", $"+strconv.Itoa(i*3+3)+")")
		args = append(args, record.UUID, record.ShortURL, record.OriginalURL)
	}

	sqlQuery := `
        INSERT INTO url_records (uuid, short_url, original_url)
        VALUES ` + strings.Join(values, ", ") + `
        ON CONFLICT (short_url) DO NOTHING
        RETURNING uuid, short_url, original_url
    `

	rows, err := tx.Query(sqlQuery, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var insertedRecord entity.URLRecord
		if err := rows.Scan(&insertedRecord.UUID, &insertedRecord.ShortURL, &insertedRecord.OriginalURL); err != nil {
			tx.Rollback()
			return nil, err
		}
		insertedRecords = append(insertedRecords, insertedRecord)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return insertedRecords, nil
}
