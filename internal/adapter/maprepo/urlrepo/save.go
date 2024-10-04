package urlrepo

import (
	"bufio"
	"encoding/json"
	"os"
	"slices"

	"github.com/google/uuid"
)

func (r *Repository) Save(id string, url string) error {
	index := slices.IndexFunc(r.records, func(record URLRecord) bool {
		return record.ShortURL == id
	})

	if index != -1 {
		return nil
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

	_, err = writer.WriteString(string(data) + "\n")
	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}
