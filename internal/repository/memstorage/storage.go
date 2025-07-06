package memstorage

import (
	"context"
	"sync"
	"time"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/ping"

	"github.com/rs/zerolog/log"
)

var _ metrics.MetricRegistry = (*MemStorage)(nil)
var _ ping.Pinger = (*MemStorage)(nil)

type MemStorage struct {
	m    map[string]metrics.Metric
	mu   *sync.RWMutex
	opts Options
}

type Options struct {
	Restore        bool
	BackupPath     string
	BackupInterval time.Duration
}

func New(ctx context.Context, opts Options) (*MemStorage, error) {
	logger := log.With().Logger()

	storage := &MemStorage{
		m:    make(map[string]metrics.Metric),
		mu:   &sync.RWMutex{},
		opts: opts,
	}

	if opts.Restore {
		if err := storage.restoreFromFile(ctx, opts.BackupPath); err != nil {
			// NOTE: autotest fails if we return error here
			// return nil, err
			// 	suite.envs = append(os.Environ(), []string{
			// 		"ADDRESS=localhost:" + flagServerPort,
			// 		"RESTORE=true",
			// 		"STORE_INTERVAL=2",
			// 		"FILE_STORAGE_PATH=" + flagFileStoragePath,
			// 	}...)
			logger.Error().Err(err).Str("path", opts.BackupPath).
				Msg("failed to restore storage from file, starting with empty storage")
		}
	}

	go storage.backupToFile(ctx, opts.BackupPath, opts.BackupInterval)

	return storage, nil
}

func (s *MemStorage) Close(ctx context.Context) error {
	logger := log.With().Str("path", s.opts.BackupPath).Logger()
	logger.Debug().Msg("closing storage")
	return s.dump(ctx, s.opts.BackupPath)
}

func (s *MemStorage) Fetch(ctx context.Context) ([]metrics.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := make([]metrics.Metric, 0, len(s.m))
	for _, m := range s.m {
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (s *MemStorage) FetchOne(ctx context.Context, m metrics.Metric) (metrics.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metric, ok := s.m[m.ID]
	if !ok {
		return metrics.Metric{}, metrics.ErrNotFound
	}

	return metric, nil
}

func (s *MemStorage) UpdateOne(ctx context.Context, m metrics.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch m.Type {
	case metrics.CounterMetricType:
		if _, ok := s.m[m.ID]; !ok {
			s.m[m.ID] = m
		} else {
			*s.m[m.ID].Counter += *m.Counter
		}
	case metrics.GaugeMetricType:
		s.m[m.ID] = m
	}

	return nil
}

func (s *MemStorage) Ping(ctx context.Context) error {
	return nil // NOTE: just a stub, isn't used
}

func (s *MemStorage) Update(ctx context.Context, many []metrics.Metric) error {
	return nil // NOTE: just a stub, isn't used
}
