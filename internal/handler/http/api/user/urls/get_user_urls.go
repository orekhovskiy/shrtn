package userurls

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type URLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (h *Handler) getUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authService.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Info("no user ID provided, rejecting")
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	urls, err := h.urlService.GetUserURLs(userID)
	if err != nil {
		h.logger.Error("unable to fetch URLs for user",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := make([]URLResponse, len(urls))
	for i, url := range urls {
		response[i] = URLResponse{
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
		}
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("failed to write JSON response",
			zap.Error(err),
		)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
