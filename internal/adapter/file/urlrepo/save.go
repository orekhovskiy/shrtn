package urlrepo

import (
	"bufio"
	"encoding/json"
	"github.com/orekhovskiy/shrtn/internal/entity"
	e "github.com/orekhovskiy/shrtn/internal/errors"
	"os"

	"github.com/google/uuid"
)

func (r *FileURLRepository) Save(id string, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.records[id]; exists {
		return &e.URLConflictError{
			ShortURL:    id,
			OriginalURL: url,
		}
	}

	record := entity.URLRecord{
		UUID:        uuid.New().String(),
		ShortURL:    id,
		OriginalURL: url,
	}

	r.records[id] = record

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
