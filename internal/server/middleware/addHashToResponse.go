package middleware

import (
	"bytes"
	"github.com/mobypolo/ya-41go/internal/shared/utils"
	"net/http"
)

func AddHashToResponse(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			hw := &hashResponseWriter{
				ResponseWriter: w,
				status:         http.StatusOK, // по умолчанию
			}

			next.ServeHTTP(hw, r)

			hash := utils.HashBody(hw.body.Bytes(), key)
			w.Header().Set("HashSHA256", hash)
		})
	}
}

type hashResponseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (w *hashResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *hashResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
