package urlrepo

import (
	"bufio"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/orekhovskiy/shrtn/internal/entity"
	"os"
)

func (r *FileURLRepository) SaveMany(records []entity.URLRecord, userID string) ([]entity.URLRecord, error) {
	if len(records) == 0 {
		return make([]entity.URLRecord, 0), nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	var responses []entity.URLRecord
	for _, record := range records {
		if _, exists := r.records[record.ShortURL]; exists {
			continue
		}
		if record.UUID == "" {
			record.UUID = uuid.New().String()
		}

		record.UserID = userID

		data, err := json.Marshal(record)
		if err != nil {
			return nil, err
		}

		if _, err := writer.WriteString(string(data) + "\n"); err != nil {
			return nil, err
		}
		responses = append(responses, record)
	}

	return responses, writer.Flush()
}
