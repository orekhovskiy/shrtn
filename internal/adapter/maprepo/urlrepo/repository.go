package urlrepo

import (
	"sync"

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
	mu       sync.RWMutex
}

func NewRepository(opts config.Config) *Repository {
	return &Repository{
		records:  []URLRecord{},
		filePath: opts.FilePath,
	}
}
