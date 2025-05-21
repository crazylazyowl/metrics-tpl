package memstorage

import "github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

var _ metrics.Repository = (*MemStorage)(nil)

type MemStorage struct {
	counters map[string][]metrics.Counter
	gagues   map[string]metrics.Gauge
}

func NewStorage() *MemStorage {
	return &MemStorage{
		counters: make(map[string][]metrics.Counter),
		gagues:   make(map[string]metrics.Gauge),
	}
}

func (s *MemStorage) UpdateCounter(name string, value metrics.Counter) error {
	// if _, ok := s.counters[name]; !ok {
	// 	s.counters[name] = make([]metrics.Counter, 0)
	// }
	s.counters[name] = append(s.counters[name], value)
	return nil
}

func (s *MemStorage) UpdateGuage(name string, value metrics.Gauge) error {
	s.gagues[name] = value
	return nil
}
