package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	ticker  *time.Ticker
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(clearInterval time.Duration) *Cache {
	cache := Cache{
		entries: map[string]cacheEntry{},
		ticker:  nil,
	}

	cache.reapLoop(clearInterval)

	return &cache
}

var mux = &sync.Mutex{}

func (cache *Cache) Add(key string, val []byte) {
	mux.Lock()
	defer mux.Unlock()

	cache.entries[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	mux.Lock()
	defer mux.Unlock()

	if cacheEntry, ok := cache.entries[key]; ok {
		return cacheEntry.val, true
	}

	return []byte{}, false
}

func (cache *Cache) reapLoop(clearInterval time.Duration) {
	if cache.ticker != nil {
		return
	}

	ticker := time.NewTicker(clearInterval)
	cache.ticker = ticker

	go func() {
		for {
			time.Sleep(clearInterval)
			t := <-ticker.C
			fmt.Printf("Current time is %s\n", t)
			mux.Lock()
			for key, entry := range cache.entries {
				if time.Since(entry.createdAt) >= clearInterval {
					fmt.Printf("Deletes key %s\n", key)
					delete(cache.entries, key)
				}
			}
			mux.Unlock()
		}
	}()
}
