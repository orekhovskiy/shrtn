package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) AddRoutes(r *chi.Mux) {
	r.Post("/", h.CreateShortURL)
	r.Get("/*", h.RedirectToOriginal)
	r.MethodNotAllowed(methodNotAllowed)
}

// Override method not allowed to reply with 400 Bad Request
func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
