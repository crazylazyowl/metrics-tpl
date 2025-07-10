package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressResponseWriter struct {
	http.ResponseWriter
	io.Writer
	compressible bool
}

func (w *compressResponseWriter) WriteHeader(status int) {
	switch w.Header().Get("Content-Type") {
	case jsonContentType, textHTMLContentType:
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		w.compressible = true
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *compressResponseWriter) Write(data []byte) (int, error) {
	if w.compressible {
		return w.Writer.Write(data)
	}
	return w.ResponseWriter.Write(data)
}

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		writer := gzip.NewWriter(w)
		defer writer.Close()
		next.ServeHTTP(&compressResponseWriter{ResponseWriter: w, Writer: writer}, r)
	})
}
