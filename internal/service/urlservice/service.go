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

type Service struct {
	urlRepository Repository
	options       config.Config
}

func NewService(opts config.Config, urlRepository Repository) *Service {
	return &Service{
		options:       opts,
		urlRepository: urlRepository,
	}
}
