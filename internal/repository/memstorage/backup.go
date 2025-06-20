package memstorage

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type snapshot struct {
	Counters map[string][]int64 `json:"counters"`
	Gauges   map[string]float64 `json:"gauges"`
}

func (s *MemStorage) Restore(path string) error {
	logger := log.With().Str("path", path).Logger()

	logger.Debug().Msg("restore from backup")

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var snapshot snapshot

	if err := json.NewDecoder(f).Decode(&snapshot); err != nil {
		return err
	}

	for key, values := range snapshot.Counters {
		for _, value := range values {
			s.counters.Append(key, value)
		}
	}

	for key, value := range snapshot.Gauges {
		s.gauges.Set(key, value)
	}

	return nil
}

func (s *MemStorage) Backup(ctx context.Context, path string, interval int) error {
	logger := log.With().Str("path", path).Int("interval", interval).Logger()

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	logger.Debug().Msg("start backup")

	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer t.Stop()

	dump := func() error {
		if err := f.Truncate(0); err != nil {
			return err
		}
		if _, err := f.Seek(0, 0); err != nil {
			return err
		}
		snapshot := snapshot{
			Counters: s.GetCounters(),
			Gauges:   s.GetGauges(),
		}
		if err := json.NewEncoder(f).Encode(&snapshot); err != nil {
			return err
		}
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			logger.Debug().Msg("stop backup")
			if err := dump(); err != nil {
				logger.Error().Err(err).Msg("faield to backup storage after context cancalation")
			}
			return ctx.Err()
		case <-t.C:
		}
		logger.Debug().Msg("backup")
		if err := dump(); err != nil {
			logger.Error().Err(err).Msg("failed to backup storage")
		}
	}
}
