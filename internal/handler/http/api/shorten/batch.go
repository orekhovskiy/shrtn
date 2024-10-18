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

	responses, err := h.urlService.ProcessBatch(batchRequests)
	if err != nil {
		h.logger.Error("failed to process batch",
			zap.Error(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
