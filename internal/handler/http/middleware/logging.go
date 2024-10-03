package middleware

import (
	"net/http"
	"time"

	"github.com/orekhovskiy/shrtn/internal/logger"

	"go.uber.org/zap"
)

func LoggingMiddleware(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapper := responseWriterWrapper{ResponseWriter: w, statusCode: 0}

			next.ServeHTTP(&wrapper, r)

			duration := time.Since(start)
			log.Info("request completed",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", wrapper.statusCode),
				zap.Int("response_size", wrapper.responseSize),
				zap.Duration("duration", duration),
			)
		})
	}
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriterWrapper) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.responseSize += size
	return size, err
}
