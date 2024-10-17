package urlservice

type Repository interface {
	Save(id string, url string) error
	GetByID(id string) (string, error)
}

type Service struct {
	urlRepository Repository
}

func NewService(urlRepository Repository) *Service {
	return &Service{
		urlRepository: urlRepository,
	}
}
