package urlrepo

import "github.com/orekhovskiy/shrtn/internal/entity"

func (r *PostgresURLRepository) GetUserURLs(userID string) ([]entity.URLRecord, error) {
	rows, err := r.db.Query(getRecordsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []entity.URLRecord
	for rows.Next() {
		var record entity.URLRecord
		if err := rows.Scan(&record.ShortURL, &record.OriginalURL); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
