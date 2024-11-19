package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	interval time.Duration
	values   map[string]cacheEntry
	mu       sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{interval: interval, values: make(map[string]cacheEntry)}
	go c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.values[key]
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, val := range c.values {
				if time.Since(val.createdAt) > c.interval {
					delete(c.values, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
