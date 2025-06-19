package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type CompressResponseWriter struct {
	http.ResponseWriter
	io.Writer
	compressible bool
}

func (w *CompressResponseWriter) WriteHeader(status int) {
	switch w.Header().Get("Content-Type") {
	case jsonContentType, textHtmlContentType:
		w.compressible = true
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *CompressResponseWriter) Write(data []byte) (int, error) {
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
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(&CompressResponseWriter{ResponseWriter: w, Writer: writer}, r)
	})
}
