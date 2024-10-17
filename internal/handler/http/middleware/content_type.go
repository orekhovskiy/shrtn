package middleware

import (
	"mime"
	"net/http"
)

func ContentTypeMiddleware(allowedContentTypes []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestContentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			isAllowed := false
			for _, contentType := range allowedContentTypes {
				if requestContentType == contentType {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
