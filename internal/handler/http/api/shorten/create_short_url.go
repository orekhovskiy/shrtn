package shorten

import (
	"encoding/json"
	"fmt"
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

	var req ShortenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(req.URL)
	if _, err = url.ParseRequestURI(originalURL); err != nil {
		h.logger.Error("unable to shorten non-url like string",
			zap.String("url", originalURL),
			zap.Error(err),
		)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := h.urlService.Save(originalURL)
	if err != nil {
		h.logger.Error("error while saving url",
			zap.String("url", originalURL),
			zap.Error(err),
		)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
	shortURL := fmt.Sprintf("%s/%s", h.opts.BaseURL, id)

	response := ShortenResponse{Result: shortURL}
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}
