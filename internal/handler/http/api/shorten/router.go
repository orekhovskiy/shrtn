package shorten

import (
	"github.com/go-chi/chi/v5"

	"github.com/orekhovskiy/shrtn/internal/handler/http/middleware"
)

func (h *Handler) AddRoutes(r *chi.Mux) {
	r.
		With(middleware.ContentTypeMiddleware([]string{
			ContentTypeJSON,
			ContentTypeGzip,
		})).
		Post("/api/shorten", h.CreateShortURL)
}
