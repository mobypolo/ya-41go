package middleware

import (
	"bytes"
	"crypto/hmac"
	"github.com/mobypolo/ya-41go/internal/shared/utils"
	"io"
	"log"
	"net/http"
)

func HashSHA256(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key != "" {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "read error", http.StatusInternalServerError)
					return
				}
				defer func() {
					if err := r.Body.Close(); err != nil {
						log.Printf("error closing response body: %v", err)
					}
				}()
				r.Body = io.NopCloser(bytes.NewReader(body))

				expectedHash := r.Header.Get("HashSHA256")
				actualHash := utils.HashBody(body, key)

				if !hmac.Equal([]byte(expectedHash), []byte(actualHash)) {
					http.Error(w, "invalid hash", http.StatusBadRequest)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
