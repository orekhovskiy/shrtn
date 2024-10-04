package api

import (
	"fmt"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), ContentTypePlainText) &&
		r.Header.Get("Content-Type") != ContentTypeGzip {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

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

	id := h.urlService.Save(originalURL)
	shortURL := fmt.Sprintf("%s/%s", h.opts.BaseURL, id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", ContentTypePlainText)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}
