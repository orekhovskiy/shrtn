package api

import (
	"go.uber.org/zap"
	"net/http"
)

func (h Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.urlService.Ping(); err != nil {
		h.logger.Error("no database connection",
			zap.Error(err),
		)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
