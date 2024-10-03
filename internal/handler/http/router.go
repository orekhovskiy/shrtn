package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api"
	"github.com/orekhovskiy/shrtn/internal/logger"
)

type Router struct {
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{router: chi.NewMux()}
}

func (r *Router) WithHandler(logger logger.Logger, h api.Handler) *Router {
	h.AddRoutes(logger, r.router)
	return r
}
