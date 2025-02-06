package utils

import "sync"

type MutexCounter struct {
	mu    sync.Mutex
	value int
}

func (c *MutexCounter) Get() (cnt int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cnt = c.value
	return
}

func (c *MutexCounter) Up() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

func (c *MutexCounter) Down() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value--
}

func (c *MutexCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}
