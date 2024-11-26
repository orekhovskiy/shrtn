package userurls

import (
	"encoding/json"
	"github.com/orekhovskiy/shrtn/internal/service/authservice"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) deleteUrls(w http.ResponseWriter, r *http.Request) {
	// Authorization check
	userID, ok := r.Context().Value(authservice.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var shortURLs []string
	if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Mark URLs as deleted
	if err := h.urlService.MarkURLsAsDeleted(shortURLs, userID); err != nil {
		h.logger.Error("failed to mark URLs as deleted", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
