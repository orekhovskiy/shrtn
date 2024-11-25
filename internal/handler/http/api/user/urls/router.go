package userurls

import (
	"github.com/go-chi/chi/v5"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeGzip = "application/x-gzip"
)

func (h *Handler) AddRoutes(r *chi.Mux) {
	r.Get("/api/user/urls", h.getUserURLs)
}
