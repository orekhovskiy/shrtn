package urlrepo

import (
	"sync"

	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/entity"
)

type Repository struct {
	records  map[string]entity.URLRecord
	filePath string
	mu       sync.RWMutex
}

func NewRepository(opts config.Config) *Repository {
	return &Repository{
		records:  make(map[string]entity.URLRecord),
		filePath: opts.FilePath,
	}
}
