package userurls

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) getUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authService.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Info("no user ID provided, rejecting")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urls, err := h.urlService.GetUserURLs(userID)
	if err != nil {
		h.logger.Error("Unable to fetch URLs for user",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := make([]map[string]string, len(urls))
	for i, url := range urls {
		response[i] = map[string]string{
			"short_url":    url.ShortURL,
			"original_url": url.OriginalURL,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
