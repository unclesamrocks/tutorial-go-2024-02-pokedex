package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	cache  map[string]cacheEntry
	ticker *time.Ticker
	mux    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(clearInterval time.Duration) *Cache {
	c := Cache{
		cache:  map[string]cacheEntry{},
		ticker: nil,
		mux:    &sync.Mutex{},
	}

	go c.reapLoop(clearInterval)

	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	val, ok := c.cache[key]
	return val.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	c.ticker = ticker
	for {
		time.Sleep(interval)
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(now time.Time, last time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	t := <-c.ticker.C
	fmt.Printf("\n[reap] %s\n", t)
	for k, v := range c.cache {
		if v.createdAt.Before(now.Add(-last)) {
			fmt.Printf("[reap][delete] %s\n", k)
			delete(c.cache, k)
		}
	}
}
