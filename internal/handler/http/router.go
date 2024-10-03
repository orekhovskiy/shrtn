package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api"
)

type Router struct {
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{router: chi.NewMux()}
}

func (r *Router) WithHandler(h api.Handler) *Router {
	h.AddRoutes(r.router)

	return r
}
