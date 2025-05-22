package memstorage

import (
	"maps"
	"slices"
	"sync"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

type counters struct {
	m  map[string][]metrics.Counter
	mu sync.RWMutex
}

func (c *counters) Copy() map[string][]metrics.Counter {
	c.mu.Lock()
	defer c.mu.Unlock()
	return maps.Clone(c.m) // TODO: shallow copy
}

func (c *counters) Get(name string) []metrics.Counter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.m[name]; !ok {
		return nil
	}
	return slices.Clone(c.m[name])
}

func (c *counters) Append(name string, value metrics.Counter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[name] = append(c.m[name], value)
}
