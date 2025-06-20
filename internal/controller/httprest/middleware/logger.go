package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type loggerResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *loggerResponseWriter) Write(data []byte) (int, error) {
	rw.size += len(data)
	return rw.ResponseWriter.Write(data)
}

func (rw *loggerResponseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		rw := loggerResponseWriter{ResponseWriter: w}
		next.ServeHTTP(&rw, r)

		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", rw.status).
			Int("size", rw.size).
			Dur("duration", time.Since(t)).
			Send()
	})
}
