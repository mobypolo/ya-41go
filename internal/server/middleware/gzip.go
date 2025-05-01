package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipDecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "failed to create gzip reader", http.StatusBadRequest)
				return
			}
			defer func(gz *gzip.Reader) {
				err := gz.Close()
				if err != nil {
					http.Error(w, "failed to close gzip reader", http.StatusBadRequest)
				}
			}(gz)
			r.Body = io.NopCloser(gz)
		}
		next.ServeHTTP(w, r)
	})
}

func GzipCompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")

		gz := gzip.NewWriter(w)
		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				http.Error(w, "failed to close gzip reader", http.StatusBadRequest)
			}
		}(gz)

		gw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gw, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
