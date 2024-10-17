package urlrepo

import (
	"bufio"
	"encoding/json"
	"os"
)

func (r *Repository) LoadAll() error {
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

	records := make(map[string]URLRecord)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record URLRecord
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return err
		}
		records[record.ShortURL] = record
	}
	r.records = records
	return scanner.Err()
}
