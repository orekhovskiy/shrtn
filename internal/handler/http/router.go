package http

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	router *chi.Mux
}

type Handler interface {
	AddRoutes(r *chi.Mux)
}

func NewRouter() *Router {
	return &Router{router: chi.NewMux()}
}

func (r *Router) WithHandler(h Handler) *Router {
	h.AddRoutes(r.router)
	return r
}
