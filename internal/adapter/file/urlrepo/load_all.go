package urlrepo

import (
	"bufio"
	"encoding/json"
	"github.com/orekhovskiy/shrtn/internal/entity"
	"os"
)

func (r *FileURLRepository) LoadAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	records := make(map[string]entity.URLRecord)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record entity.URLRecord
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return err
		}
		records[record.ShortURL] = record
	}
	r.records = records
	return scanner.Err()
}
