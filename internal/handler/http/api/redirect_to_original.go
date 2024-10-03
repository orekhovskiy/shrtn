package api

import (
	"log"
	"net/http"
	"strings"
)

func (h Handler) RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, err := h.urlService.GetByID(id)

	if err != nil {
		log.Printf("error getting original url by id %s: %v", id, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Printf("Redirecting to %s", originalURL)
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
