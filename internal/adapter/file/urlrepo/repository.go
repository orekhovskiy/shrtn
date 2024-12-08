package urlrepo

import (
	"sync"

	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/entity"
)

type FileURLRepository struct {
	records  map[string]entity.URLRecord
	filePath string
	mu       sync.RWMutex
}

func NewRepository(opts config.Config) *FileURLRepository {
	return &FileURLRepository{
		records:  make(map[string]entity.URLRecord),
		filePath: opts.FilePath,
	}
}
