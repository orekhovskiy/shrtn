package api

import (
	"fmt"
	"net/http"
	"strings"
)

func (h Handler) RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, err := h.urlService.GetById(id)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	logMessage := fmt.Sprintf("Redirecting to %s", originalURL)
	fmt.Println(logMessage)
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
