package memstorage

import (
	"maps"
	"sync"
)

type counters struct {
	m  map[string]int64
	mu *sync.RWMutex
}

func newCounters() *counters {
	return &counters{
		m:  make(map[string]int64),
		mu: &sync.RWMutex{},
	}
}

func (c *counters) Copy() map[string]int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return maps.Clone(c.m)
}

func (c *counters) Get(name string) (int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.m[name]
	return value, ok
}

func (c *counters) Update(name string, value int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	tmp := c.m[name]
	c.m[name] = tmp + value

	return c.m[name]
}
