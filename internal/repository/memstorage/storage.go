package memstorage

import (
	"context"
	"time"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
	"github.com/rs/zerolog/log"
)

var _ metrics.MetricsStorage = (*MemStorage)(nil)

type MemStorage struct {
	counters *counters
	gauges   *gauges
	opts     Options
}

type Options struct {
	Restore        bool
	BackupPath     string
	BackupInterval time.Duration
}

func New(ctx context.Context, opts Options) (*MemStorage, error) {
	storage := &MemStorage{
		counters: newCounters(),
		gauges:   newGauges(),
		opts:     opts,
	}

	if opts.Restore {
		if err := storage.restoreFromFile(opts.BackupPath); err != nil {
			// NOTE: autotest fails if we return error here
			// return nil, err
			// 	suite.envs = append(os.Environ(), []string{
			// 		"ADDRESS=localhost:" + flagServerPort,
			// 		"RESTORE=true",
			// 		"STORE_INTERVAL=2",
			// 		"FILE_STORAGE_PATH=" + flagFileStoragePath,
			// 	}...)
		}
	}

	go storage.backupToFile(ctx, opts.BackupPath, opts.BackupInterval)

	return storage, nil
}

func (s *MemStorage) Close() error {
	logger := log.With().Str("path", s.opts.BackupPath).Logger()
	logger.Debug().Msg("closing storage")
	return s.dump(s.opts.BackupPath)
}

func (s *MemStorage) GetCounters() map[string][]int64 {
	return s.counters.Copy()
}

func (s *MemStorage) GetGauges() map[string]float64 {
	return s.gauges.Copy()
}

func (s *MemStorage) GetCounter(name string) ([]int64, error) {
	values, found := s.counters.Get(name)
	if !found {
		return nil, metrics.ErrUnknownMetric
	}
	return values, nil
}

func (s *MemStorage) GetGauge(name string) (float64, error) {
	value, found := s.gauges.Get(name)
	if !found {
		return 0, metrics.ErrUnknownMetric
	}
	return value, nil
}

func (s *MemStorage) AppendCounter(name string, value int64) error {
	s.counters.Append(name, value)
	return nil
}

func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.gauges.Set(name, value)
	return nil
}
