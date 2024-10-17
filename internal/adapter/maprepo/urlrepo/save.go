package urlrepo

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

func (r *Repository) Save(id string, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, record := range r.records {
		if record.ShortURL == id {
			return nil
		}
	}

	record := URLRecord{
		UUID:        uuid.New().String(),
		ShortURL:    id,
		OriginalURL: url,
	}

	r.records = append(r.records, record)

	if r.filePath == "" {
		return nil
	}

	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	if _, err := writer.WriteString(string(data) + "\n"); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
