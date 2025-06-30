package memstorage

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *MemStorage) restoreFromFile(ctx context.Context, path string) error {
	logger := log.With().Str("path", path).Logger()

	logger.Debug().Msg("restore from backup")

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := json.NewDecoder(f).Decode(&s.m); err != nil {
		return err
	}

	return nil
}

func (s *MemStorage) backupToFile(ctx context.Context, path string, dur time.Duration) error {
	logger := log.With().
		Str("path", path).
		Dur("dur", dur).
		Logger()

	logger.Debug().Msg("start backup")

	t := time.NewTicker(dur)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Debug().Msg("stop backup")
			return ctx.Err()
		case <-t.C:
			logger.Debug().Msg("backup")
			if err := s.dump(ctx, path); err != nil {
				logger.Error().Err(err).Msg("failed to backup storage")
			}
		}
	}
}

func (s *MemStorage) dump(ctx context.Context, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := json.NewEncoder(f).Encode(&s.m); err != nil {
		return err
	}

	return nil
}
