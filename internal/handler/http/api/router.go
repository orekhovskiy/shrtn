package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/orekhovskiy/shrtn/internal/handler/http/middleware"
)

func (h *Handler) AddRoutes(r *chi.Mux) {
	r.Use(middleware.LoggingMiddleware(h.logger))
	r.Use(middleware.GzipMiddleware(h.logger))

	r.Get("/ping", h.Ping)
	r.
		With(middleware.AuthMiddleware(h.opts, h.logger, false)).
		With(middleware.ContentTypeMiddleware([]string{
			ContentTypePlainText,
			ContentTypeGzip,
		})).
		Post("/", h.CreateShortURL)
	r.Get("/*", h.RedirectToOriginal)
	r.MethodNotAllowed(methodNotAllowed)
}

// Override method not allowed to reply with 400 Bad Request
func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
