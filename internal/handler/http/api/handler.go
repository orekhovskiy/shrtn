package api

import (
	"github.com/orekhovskiy/shrtn/config"
	"go.uber.org/zap"
)

type Service interface {
	GetByID(id string) (url string, err error)
	Save(url string) (id string)
}

type Handler struct {
	opts       config.Config
	urlService Service
	logger     zap.Logger
}

func NewHandler(logger *zap.Logger, opts *config.Config, urlService Service) *Handler {
	return &Handler{logger: *logger, opts: *opts, urlService: urlService}
}
