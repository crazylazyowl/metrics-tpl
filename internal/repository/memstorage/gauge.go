package memstorage

import (
	"maps"
	"sync"
)

type gauges struct {
	m  map[string]float64
	mu *sync.RWMutex
}

func (g *gauges) Copy() map[string]float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return maps.Clone(g.m)
}

func (g *gauges) Get(name string) (float64, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	value, ok := g.m[name]
	return value, ok
}

func (g *gauges) Set(name string, value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.m[name] = value
}
