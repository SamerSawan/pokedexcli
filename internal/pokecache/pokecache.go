package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	newEntry := cacheEntry{createdAt: time.Now(), val: val}
	c.entries[key] = newEntry
	c.mu.Unlock()
	return
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, exists := c.entries[key]
	c.mu.Unlock()
	if exists {
		return entry.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)

	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, val := range c.entries {
			if time.Since(val.createdAt) > interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}
