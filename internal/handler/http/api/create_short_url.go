package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (h Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
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
		errorMessage := fmt.Sprintf("unable to shorten non-url like string %s: %s", originalURL, err)
		fmt.Println(errorMessage)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id := h.urlService.Save(originalURL)
	shortURL := fmt.Sprintf("%s/%s", h.opts.BaseURL, id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(shortURL))
}
