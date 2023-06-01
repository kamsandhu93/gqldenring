package memDB

import "sync"

type counter struct {
	count int
	mu    sync.RWMutex
}

func newCounter(start int) *counter {
	return &counter{
		count: start,
		mu:    sync.RWMutex{},
	}
}

func (c *counter) read() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}

func (c *counter) increment() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count += 1
	return c.count
}
