package middleware

import (
	"net/http"
	"time"

	"log"
)

type LoggerResponseWriter struct {
	http.ResponseWriter
	custom struct {
		status int
		size   int
	}
}

func (rw *LoggerResponseWriter) Write(data []byte) (int, error) {
	rw.custom.size += len(data)
	return rw.ResponseWriter.Write(data)
}

func (rw *LoggerResponseWriter) WriteHeader(statusCode int) {
	rw.custom.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		rw := LoggerResponseWriter{ResponseWriter: w}
		next.ServeHTTP(&rw, r)
		log.Printf("method=%s url=%s status=%d size=%d duration=%v\n",
			r.Method, r.URL.String(), rw.custom.status, rw.custom.size, time.Since(t))
	})
}
