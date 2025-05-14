package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/mobypolo/ya-41go/internal/shared/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
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

		var bodyPreview any

		if r.Body != nil && r.ContentLength != 0 {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
					gr, err := gzip.NewReader(bytes.NewReader(bodyBytes))
					if err == nil {
						bodyBytes, _ = io.ReadAll(gr)
						_ = gr.Close()
					}
				}

				var parsed any
				if err := json.Unmarshal(bodyBytes, &parsed); err == nil {
					bodyPreview = parsed // structured object
				} else {
					bodyPreview = string(bodyBytes) // fallback: raw string
				}

				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		logger.L().Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", rw.statusCode),
			zap.Int("size", rw.size),
			zap.Duration("duration", duration),
			zap.Any("body", bodyPreview),
		)
	})
}
