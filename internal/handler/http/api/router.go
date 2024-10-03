package api

import (
	"net/http"

	"github.com/orekhovskiy/shrtn/internal/handler/http/middleware"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) AddRoutes(logger *zap.Logger, r *chi.Mux) {
	r.Use(middleware.LoggingMiddleware(logger))

	r.Post("/", h.CreateShortURL)
	r.Get("/*", h.RedirectToOriginal)
	r.MethodNotAllowed(methodNotAllowed)
}

// Override method not allowed to reply with 400 Bad Request
func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
