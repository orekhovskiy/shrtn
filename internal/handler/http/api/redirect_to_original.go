package api

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func (h Handler) RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, err := h.urlService.GetByID(id)

	if err != nil {
		h.logger.Error("error getting original url by id",
			zap.String("id", id),
			zap.Error(err),
		)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	h.logger.Info("redirecting",
		zap.String("url", originalURL),
	)
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
