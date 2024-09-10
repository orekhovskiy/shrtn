package api

import (
	service "github.com/orekhovskiy/shrtn/internal/service/urlservice"
)

type Handler struct {
	urlService service.Service
}

func NewHandler(urlService service.Service) *Handler {
	return &Handler{urlService: urlService}
}
