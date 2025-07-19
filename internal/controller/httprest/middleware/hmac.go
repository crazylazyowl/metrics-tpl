package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/security"

	"github.com/rs/zerolog/log"
)

func CheckHMAC(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if have := r.Header.Get("HashSHA256"); have != "" {
				logger := log.With().Logger()
				data, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Error().Err(err).Msg("failed to read request body")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				want := security.HMACString([]byte(secret), data)
				if have != want {
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
