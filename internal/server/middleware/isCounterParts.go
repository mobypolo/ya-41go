package middleware

import (
	"github.com/mobypolo/ya-41go/internal/server/helpers"
	"net/http"
)

func RequirePathParts(minParts int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверка пути
			parts := helpers.SplitStrToChunks(r.URL.Path)
			if len(parts) < minParts {
				http.Error(w, "invalid path", http.StatusNotFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
