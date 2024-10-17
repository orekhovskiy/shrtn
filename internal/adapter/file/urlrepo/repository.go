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
	records  map[string]URLRecord
	filePath string
	mu       sync.RWMutex
}

func NewRepository(opts config.Config) *Repository {
	return &Repository{
		records:  make(map[string]URLRecord),
		filePath: opts.FilePath,
	}
}
