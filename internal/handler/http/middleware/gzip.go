package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	ContentTypeHTML = "text/html"
	ContentTypeJSON = "application/json"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			defer reader.Close()
			r.Body = reader
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != ContentTypeHTML && contentType != ContentTypeJSON {
			next.ServeHTTP(w, r)
			return
		}

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		gzipResponseWriter := &GzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipResponseWriter, r)
	})
}

type GzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g *GzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}
