package cache

import "sync"

type ConcurrentCache struct {
	mx      sync.Mutex
	cache   *Cache
	maxSize int64
}

func (conCache *ConcurrentCache) Put(key string, value Value) {
	conCache.mx.Lock()
	defer conCache.mx.Unlock()

	if conCache.cache == nil {
		conCache.cache = Init(conCache.maxSize)
	}
	conCache.cache.Put(key, value)
}

func (conCache *ConcurrentCache) Get(key string) (value Value, ok bool) {
	conCache.mx.Lock()
	defer conCache.mx.Unlock()

	if conCache.cache == nil {
		return
	}

	if value, ok := conCache.cache.Get(key); ok {
		return value, true
	}
	return
}
