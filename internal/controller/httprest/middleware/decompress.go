package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "failed to decompress request body", http.StatusBadRequest)
			return
		}
		defer reader.Close()
		r.Body = reader
		next.ServeHTTP(w, r)
	})
}
