package urlrepo

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/orekhovskiy/shrtn/config"
)

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Repository struct {
	records  []URLRecord
	filePath string
}

func loadAll(filePath string) ([]URLRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []URLRecord{}, nil
		}
		return []URLRecord{}, err
	}
	defer file.Close()

	var records []URLRecord
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record URLRecord
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return []URLRecord{}, err
		}
		records = append(records, record)
	}
	return records, scanner.Err()
}

func NewRepository(opts config.Config) (*Repository, error) {
	urlMapping, err := loadAll(opts.FilePath)
	if err != nil {
		return nil, err
	}
	return &Repository{
		records:  urlMapping,
		filePath: opts.FilePath,
	}, nil
}
