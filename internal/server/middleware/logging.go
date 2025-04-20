package middleware

import (
	"github.com/mobypolo/ya-41go/internal/shared/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		logger.L().Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", rw.statusCode),
			zap.Int("size", rw.size),
			zap.Duration("duration", duration),
		)
	})
}
