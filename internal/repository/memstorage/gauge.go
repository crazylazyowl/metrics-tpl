package memstorage

import (
	"maps"
	"sync"
)

type gauges struct {
	m  map[string]float64
	mu sync.RWMutex
}

func (g *gauges) Copy() map[string]float64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	return maps.Clone(g.m)
}

func (g *gauges) Get(name string) float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, ok := g.m[name]; !ok {
		return 0 // TODO: 0 is a valid metric value
	}
	return g.m[name]
}

func (g *gauges) Set(name string, value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.m[name] = value
}
