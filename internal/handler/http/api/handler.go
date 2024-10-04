package api

import (
	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/logger"
)

type Service interface {
	GetByID(id string) (url string, err error)
	Save(url string) (id string, err error)
}

type Handler struct {
	opts       config.Config
	urlService Service
	logger     logger.Logger
}

func NewHandler(logger logger.Logger, opts *config.Config, urlService Service) *Handler {
	return &Handler{logger: logger, opts: *opts, urlService: urlService}
}
