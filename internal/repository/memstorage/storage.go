package memstorage

import (
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

var _ metrics.MetricsStorage = (*MemStorage)(nil)

type MemStorage struct {
	counters counters
	gauges   gauges
}

func New() *MemStorage {
	return &MemStorage{
		counters: counters{m: make(map[string][]metrics.Counter)},
		gauges:   gauges{m: make(map[string]metrics.Gauge)},
	}
}

func (s *MemStorage) GetCounters() map[string][]metrics.Counter {
	return s.counters.Copy()
}

func (s *MemStorage) GetGauges() map[string]metrics.Gauge {
	return s.gauges.Copy()
}

func (s *MemStorage) GetCounter(name string) ([]metrics.Counter, error) {
	values := s.counters.Get(name)
	if values == nil {
		return nil, metrics.ErrUnknownMetric
	}
	return values, nil
}

func (s *MemStorage) GetGauge(name string) (metrics.Gauge, error) {
	value := s.gauges.Get(name)
	if value == 0 {
		return 0, metrics.ErrUnknownMetric
	}
	return value, nil
}

func (s *MemStorage) UpdateCounter(name string, value metrics.Counter) error {
	s.counters.Append(name, value)
	return nil
}

func (s *MemStorage) UpdateGauge(name string, value metrics.Gauge) error {
	s.gauges.Set(name, value)
	return nil
}
