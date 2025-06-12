package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
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

		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", rw.custom.status).
			Int("size", rw.custom.size).
			Dur("duration", time.Since(t)).
			Send()
	})
}
