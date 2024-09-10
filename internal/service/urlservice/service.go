package urlservice

import (
	repo "github.com/orekhovskiy/shrtn/internal/adapter/maprepo/urlrepo"
)

type Service struct {
	urlRepository repo.Repository
}

func NewService(urlRepository repo.Repository) *Service {
	return &Service{
		urlRepository: urlRepository,
	}
}
