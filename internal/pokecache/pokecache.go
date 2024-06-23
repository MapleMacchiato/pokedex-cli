package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

var cache = &Cache{}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		reapDate := time.Now().Add(-interval)
		for k, v := range c.entries {
			if v.createdAt.Compare(reapDate) <= 0 {
				delete(c.entries, k)
				fmt.Printf("\nDeleted %s from cache\n", k)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = entry
	fmt.Printf("\nAdded %s to cache\n", key)
}

func (c *Cache) get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	fmt.Printf("\nRetrieving key %s from cache. Was present: %t\n", key, ok)
	return entry.val, ok

}

func Add(key string, val []byte) {
	cache.add(key, val)
}

func Get(key string) ([]byte, bool) {
	return cache.get(key)
}
func NewCache(interval time.Duration) {
	cache.entries = make(map[string]cacheEntry)
	go cache.reapLoop(interval)
}
