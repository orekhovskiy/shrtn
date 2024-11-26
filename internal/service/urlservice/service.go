package urlservice

import (
	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/entity"
)

type Repository interface {
	Save(id string, url string, userID string) error
	GetByID(id string) (*entity.URLRecord, error)
	Ping() error
	SaveMany(recordsToSave []entity.URLRecord, userID string) ([]entity.URLRecord, error)
	GetUserURLs(userID string) ([]entity.URLRecord, error)
	MarkURLsAsDeleted(batch []string, userID string) error
}

type URLShortenerService struct {
	urlRepository Repository
	options       config.Config
}

func NewService(opts config.Config, urlRepository Repository) *URLShortenerService {
	return &URLShortenerService{
		options:       opts,
		urlRepository: urlRepository,
	}
}
