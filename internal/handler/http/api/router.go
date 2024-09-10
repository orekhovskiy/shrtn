package api

import (
	"fmt"
	"net/http"
)

func (h *Handler) AddRoutes(r *http.ServeMux) {
	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				fmt.Println("yeah get")
				h.RedirectToOriginal(w, r)
			case http.MethodPost:
				fmt.Println("yeah post")
				h.CreateShortUrl(w, r)
			default:
				http.Error(w, "Bad Request", http.StatusBadRequest)

			}
		},
	)
}
