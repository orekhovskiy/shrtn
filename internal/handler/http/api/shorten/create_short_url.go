package shorten

import (
	"encoding/json"
	e "github.com/orekhovskiy/shrtn/internal/errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeGzip = "application/x-gzip"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (h Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req ShortenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(req.URL)
	if _, err = url.ParseRequestURI(originalURL); err != nil {
		h.logger.Error("unable to shorten non-URL like string",
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
			response := ShortenResponse{Result: shortURL}
			w.Header().Set("Content-Type", ContentTypeJSON)
			w.WriteHeader(http.StatusConflict)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Internal Error", http.StatusInternalServerError)
			}
			return
		}

		h.logger.Error("error while saving URL",
			zap.String("url", originalURL),
			zap.Error(err),
		)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	shortURL := h.urlService.BuildURL(id)
	response := ShortenResponse{Result: shortURL}
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
