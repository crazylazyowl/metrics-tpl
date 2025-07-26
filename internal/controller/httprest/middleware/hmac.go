package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/security"

	"github.com/rs/zerolog/log"
)

func CheckHMAC(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if hmac := r.Header.Get("HashSHA256"); hmac != "" {
				logger := log.With().Logger()
				data, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Error().Err(err).Msg("failed to read request body")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if hmac != security.HMACString([]byte(secret), data) {
					logger.Error().Err(err).Msg("hash is invalid")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				logger.Debug().Msg("restore request body")
				r.Body = io.NopCloser(bytes.NewReader(data))
			}
			next.ServeHTTP(w, r)
		})
	}
}

type hmacResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        *bytes.Buffer
}

func (w *hmacResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *hmacResponseWriter) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

func AddHMAC(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.With().Logger()
			rw := hmacResponseWriter{ResponseWriter: w, buf: bytes.NewBuffer(nil)}
			next.ServeHTTP(&rw, r)
			hmac := security.HMACString([]byte(secret), rw.buf.Bytes())
			w.Header().Set("HashSHA256", hmac)
			w.WriteHeader(rw.statusCode)
			_, err := w.Write(rw.buf.Bytes())
			if err != nil {
				logger.Error().Err(err).Msg("failed to write response body")
			}
		})
	}
}
