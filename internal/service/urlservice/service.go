package urlservice

import "github.com/orekhovskiy/shrtn/internal/entity"

type Repository interface {
	Save(id string, url string) error
	GetByID(id string) (string, error)
	Ping() error
	SaveMany([]entity.URLRecord) ([]entity.URLRecord, error)
}

type Service struct {
	urlRepository Repository
}

func NewService(urlRepository Repository) *Service {
	return &Service{
		urlRepository: urlRepository,
	}
}
