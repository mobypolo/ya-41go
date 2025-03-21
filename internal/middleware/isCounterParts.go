package middleware

import (
	"github.com/mobypolo/ya-41go/internal/helpers"
	"net/http"
)

func RequirePathParts(minParts int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := helpers.SplitStrToChunks(r.URL.Path)

		if len(parts) < minParts {
			http.Error(w, "not enough path parts", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
