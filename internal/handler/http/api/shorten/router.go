package shorten

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) AddRoutes(r *chi.Mux) {
	r.Post("/api/shorten", h.CreateShortURL)
}
