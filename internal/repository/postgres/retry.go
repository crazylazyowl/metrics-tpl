package postgres

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

func retryOnConnectionError(ctx context.Context, f func(ctx context.Context) error) error {
	logger := log.With().Logger()
	var err error
	delay := 1
	for range 4 {
		if err = f(ctx); err == nil {
			break
		}
		if !isConnectionError(err) {
			return err
		}
		logger.Warn().Err(err).Msgf("request to DB failed; retry in %d seconds...", delay)
		select {
		case <-time.After(time.Duration(delay) * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
		delay += 2
	}
	return err
}

func isConnectionError(err error) bool {
	var connErr *net.OpError
	return errors.As(err, &connErr)
}
