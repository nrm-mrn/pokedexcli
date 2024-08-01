package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	content   map[string]cacheEntry
	mu        *sync.Mutex
	staleTime time.Duration
}

func NewCache(interval time.Duration) Cache {
	cacheInst := Cache{
		content:   make(map[string]cacheEntry),
		mu:        &sync.Mutex{},
		staleTime: interval,
	}
	go cacheInst.reapLoop()
	return cacheInst
}

func (c Cache) Add(key string, val []byte) {
	newEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.content[key] = newEntry
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Printf("Checking cache for %v\n", key)
	val, ok := c.content[key]
	if !ok {
		fmt.Printf("%v not found in cache\n", key)
		return []byte{}, false
	}
	fmt.Printf("%v found in cache\n", key)
	return val.val, true
}

func (c Cache) reapLoop() {
	ticker := time.NewTicker(c.staleTime)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mu.Lock()
		for k, v := range c.content {
			age := time.Now().Sub(v.createdAt)
			if age > c.staleTime {
				delete(c.content, k)
			}
		}
		c.mu.Unlock()
	}
}
