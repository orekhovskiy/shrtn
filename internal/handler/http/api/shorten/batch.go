package shorten

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"

	"github.com/orekhovskiy/shrtn/internal/entity"
)

func (h *Handler) Batch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var batchRequests []entity.BatchRequest
	if err := json.Unmarshal(body, &batchRequests); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if len(batchRequests) == 0 {
		http.Error(w, "Empty batch", http.StatusBadRequest)
		return
	}

	userID, ok := h.authService.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Info("no user ID provided, rejecting")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	responses, err := h.urlService.ProcessBatch(batchRequests, userID)
	if err != nil {
		h.logger.Error("failed to process batch",
			zap.Error(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responses)
}
