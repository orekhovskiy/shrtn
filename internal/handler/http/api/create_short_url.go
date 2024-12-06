package api

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"

	e "github.com/orekhovskiy/shrtn/internal/errors"
)

const (
	ContentTypePlainText = "text/plain"
	ContentTypeGzip      = "application/x-gzip"
)

func (h Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// Read the body from the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Trim spaces and validate if the original URL is a valid URI
	originalURL := strings.TrimSpace(string(body))
	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		h.logger.Error("invalid URL format", zap.String("url", originalURL), zap.Error(err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Retrieve the user ID from the context
	userID, ok := h.authService.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Info("user ID not provided, rejecting request")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Try to save the original URL and handle potential conflicts
	id, err := h.urlService.Save(originalURL, userID)
	if err != nil {
		var urlConflictError *e.URLConflictError
		if errors.As(err, &urlConflictError) {
			// Handle URL conflict by returning the existing short URL
			shortURL, err := h.urlService.BuildURL(urlConflictError.ShortURL)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", ContentTypePlainText)
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(shortURL))
			if err != nil {
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}
			return
		}

		// If the error is something else, log it and return a generic internal error
		h.logger.Error("error while saving URL", zap.String("url", originalURL), zap.Error(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	// Build the short URL from the generated ID
	shortURL, err := h.urlService.BuildURL(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send the response with the created short URL
	w.Header().Set("Content-Type", ContentTypePlainText)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, "Internal Server Errore", http.StatusInternalServerError)
	}
}
