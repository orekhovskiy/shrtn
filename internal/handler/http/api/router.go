package api

import (
	"net/http"
)

func (h *Handler) AddRoutes(r *http.ServeMux) {
	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.RedirectToOriginal(w, r)
			case http.MethodPost:
				h.CreateShortUrl(w, r)
			default:
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		},
	)
}
