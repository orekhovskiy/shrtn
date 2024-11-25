package api

import (
	e "github.com/orekhovskiy/shrtn/internal/errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

const (
	ContentTypePlainText = "text/plain"
	ContentTypeGzip      = "application/x-gzip"
)

func (h Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))

	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		h.logger.Error("unable to shorten non-url like string",
			zap.String("url", originalURL),
			zap.Error(err),
		)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID, ok := h.authService.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Info("no user ID provided, rejecting")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := h.urlService.Save(originalURL, userID)
	if err != nil {
		if urlConflictError, ok := err.(*e.URLConflictError); ok {
			shortURL := h.urlService.BuildURL(urlConflictError.ShortURL)
			w.Header().Set("Content-Type", ContentTypePlainText)
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(shortURL))
			if err != nil {
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}
			return
		}
		h.logger.Error("error while saving url",
			zap.String("url", originalURL),
			zap.Error(err),
		)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	shortURL := h.urlService.BuildURL(id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", ContentTypePlainText)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}
