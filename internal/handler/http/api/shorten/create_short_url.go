package shorten

import (
	"encoding/json"
	"errors"
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
		h.respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	defer r.Body.Close()

	var req ShortenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	originalURL := strings.TrimSpace(req.URL)
	if _, err = url.ParseRequestURI(originalURL); err != nil {
		h.logger.Error("unable to shorten non-URL like string",
			zap.String("url", originalURL),
			zap.Error(err),
		)
		h.respondWithError(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	userID, ok := h.authService.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Info("no user ID provided, rejecting")
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := h.urlService.Save(originalURL, userID)
	if err != nil {
		h.handleSaveError(w, originalURL, err)
		return
	}

	h.respondWithShortURL(w, id, http.StatusCreated)
}

func (h Handler) handleSaveError(w http.ResponseWriter, originalURL string, err error) {
	var urlConflictError *e.URLConflictError
	if errors.As(err, &urlConflictError) {
		shortURL, buildErr := h.urlService.BuildURL(urlConflictError.ShortURL)
		if buildErr != nil {
			h.logger.Error("failed to build URL",
				zap.String("shortURL", urlConflictError.ShortURL),
				zap.Error(buildErr),
			)
			h.respondWithError(w, http.StatusInternalServerError, "Internal Error")
			return
		}

		response := ShortenResponse{Result: shortURL}
		w.Header().Set("Content-Type", ContentTypeJSON)
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	h.logger.Error("error while saving URL",
		zap.String("url", originalURL),
		zap.Error(err),
	)
	h.respondWithError(w, http.StatusInternalServerError, "Internal Error")
}

func (h Handler) respondWithShortURL(w http.ResponseWriter, id string, statusCode int) {
	shortURL, err := h.urlService.BuildURL(id)
	if err != nil {
		h.logger.Error("failed to build URL",
			zap.String("id", id),
			zap.Error(err),
		)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	response := ShortenResponse{Result: shortURL}
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func (h Handler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}
