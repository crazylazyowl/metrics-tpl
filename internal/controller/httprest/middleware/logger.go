package middleware

import (
	"net/http"
	"time"

	"log"
)

type LoggerResponseWriter struct {
	http.ResponseWriter
	custome struct {
		status int
		size   int
	}
}

func (rw *LoggerResponseWriter) Write(data []byte) (int, error) {
	rw.custome.size += len(data)
	return rw.ResponseWriter.Write(data)
}

func (rw *LoggerResponseWriter) WriteHeader(statusCode int) {
	rw.custome.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		rw := LoggerResponseWriter{ResponseWriter: w}
		next.ServeHTTP(&rw, r)
		log.Printf("method=%s url=%s status=%d size=%d duration=%v\n",
			r.Method, r.URL.String(), rw.custome.status, rw.custome.size, time.Since(t))
	})
}
