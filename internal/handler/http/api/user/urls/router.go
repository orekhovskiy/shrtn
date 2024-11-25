package userurls

import (
	"github.com/go-chi/chi/v5"
	"github.com/orekhovskiy/shrtn/internal/handler/http/middleware"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeGzip = "application/x-gzip"
)

func (h *Handler) AddRoutes(r *chi.Mux) {
	r.
		With(middleware.AuthMiddleware(h.opts, h.logger, false)).
		Get("/api/user/urls", h.getUserURLs)
}
