package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api"
	"go.uber.org/zap"
)

type Router struct {
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{router: chi.NewMux()}
}

func (r *Router) WithHandler(logger *zap.Logger, h api.Handler) *Router {
	h.AddRoutes(logger, r.router)
	return r
}
