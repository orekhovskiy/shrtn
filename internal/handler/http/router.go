package http

import (
	"github.com/orekhovskiy/shrtn/internal/handler/http/api"
	"net/http"
)

type Router struct {
	router *http.ServeMux
}

func NewRouter() *Router {
	return &Router{router: http.NewServeMux()}
}

func (r *Router) WithHandler(h api.Handler) *Router {
	h.AddRoutes(r.router)

	return r
}
