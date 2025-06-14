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
		counters: counters{m: make(map[string][]int64)},
		gauges:   gauges{m: make(map[string]float64)},
	}
}

func (s *MemStorage) GetCounters() map[string][]int64 {
	return s.counters.Copy()
}

func (s *MemStorage) GetGauges() map[string]float64 {
	return s.gauges.Copy()
}

func (s *MemStorage) GetCounter(name string) ([]int64, error) {
	values := s.counters.Get(name)
	if values == nil {
		return nil, metrics.ErrUnknownMetric
	}
	return values, nil
}

func (s *MemStorage) GetGauge(name string) (float64, error) {
	value := s.gauges.Get(name)
	if value == 0 {
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
