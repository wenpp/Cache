package cache

import (
	"errors"
	"sync"
)

type CacheNameSpace struct {
	name  string
	cache *Cache
}

var (
	cacheNameSpaces = make(map[string]*CacheNameSpace)
	mx              sync.Mutex
)

func InitNameSapce(name string, spaceCap int64) *CacheNameSpace {
	mx.Lock()
	defer mx.Unlock()
	if GetNameSpace(name) == nil {
		namespace := &CacheNameSpace{
			name:  name,
			cache: Init(spaceCap),
		}

		cacheNameSpaces[name] = namespace
	}

	return GetNameSpace(name)
}

func GetNameSpace(name string) *CacheNameSpace {
	return cacheNameSpaces[name]
}

func (namespace *CacheNameSpace) Get(key string) (value Value, ok bool) {
	if key == "" {
		return nil, false
	}

	if value, ok := namespace.cache.Get(key); ok {
		return value, true
	}

	return nil, false
}

func (namespace *CacheNameSpace) Put(key string, value Value) error {
	mx.Lock()
	defer mx.Unlock()
	if key == "" {
		return errors.New("key is required")
	}

	namespace.cache.Put(key, value)

	return nil
}
