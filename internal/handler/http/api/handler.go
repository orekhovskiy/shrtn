package api

import "github.com/orekhovskiy/shrtn/config"

type Service interface {
	GetByID(id string) (url string, err error)
	Save(url string) (id string)
}

type Handler struct {
	opts       config.Config
	urlService Service
}

func NewHandler(opts *config.Config, urlService Service) *Handler {
	return &Handler{opts: *opts, urlService: urlService}
}
