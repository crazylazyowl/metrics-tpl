package memstorage

import (
	"slices"
	"sync"
)

type counters struct {
	m  map[string][]int64
	mu sync.RWMutex
}

func (c *counters) Copy() map[string][]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	m := make(map[string][]int64, len(c.m))
	for k, v := range c.m {
		m[k] = slices.Clone(v)
	}
	return m
}

func (c *counters) Get(name string) []int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.m[name]; !ok {
		return nil
	}
	return slices.Clone(c.m[name])
}

func (c *counters) Append(name string, value int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[name] = append(c.m[name], value)
}
