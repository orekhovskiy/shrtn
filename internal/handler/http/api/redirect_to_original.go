package api

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func (h Handler) RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	result, err := h.urlService.GetByID(id)

	if err != nil {
		h.logger.Error("error getting original url by id",
			zap.String("id", id),
			zap.Error(err),
		)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if !result.Success {
		h.logger.Info("requesting deleted url",
			zap.String("url", result.OriginalURL),
		)
		http.Error(w, "Gone", http.StatusGone)
		return
	}

	h.logger.Info("redirecting",
		zap.String("url", result.OriginalURL),
	)
	http.Redirect(w, r, result.OriginalURL, http.StatusTemporaryRedirect)
}
