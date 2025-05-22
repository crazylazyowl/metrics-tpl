package memstorage

import (
	"maps"
	"sync"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

type gauges struct {
	m  map[string]metrics.Gauge
	mu sync.RWMutex
}

func (g *gauges) Copy() map[string]metrics.Gauge {
	g.mu.Lock()
	defer g.mu.Unlock()
	return maps.Clone(g.m)
}

func (g *gauges) Get(name string) metrics.Gauge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if _, ok := g.m[name]; !ok {
		return 0 // TODO: 0 is a valid metric value
	}
	return g.m[name]
}

func (g *gauges) Set(name string, value metrics.Gauge) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.m[name] = value
}
