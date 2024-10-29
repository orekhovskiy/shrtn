package urlservice

import (
	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/entity"
)

type Repository interface {
	Save(id string, url string) error
	GetByID(id string) (string, error)
	Ping() error
	SaveMany([]entity.URLRecord) ([]entity.URLRecord, error)
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
