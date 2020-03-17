package cache

import "sync"

type Cache struct {
	mx      sync.Mutex
	core    *CoreCache
	maxSize int64
}

func (cache *Cache) Put(key string, value Value) {
	cache.mx.Lock()
	defer cache.mx.Unlock()

	if cache.core == nil {
		cache.core = Init(cache.maxSize)
	}

	cache.core.Put(key, value)
}

func (cache *Cache) Get(key string) (value Value, ok bool) {
	cache.mx.Lock()
	defer cache.mx.Unlock()

	if cache.core == nil {
		return
	}

	if value, ok := cache.core.Get(key); ok {
		return value, true
	}
	return
}
