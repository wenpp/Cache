package cache

import (
	"errors"
	"fmt"
	"sync"
)

type CacheNameSpace struct {
	name       string     //缓存控件名称
	cache      *Cache     //当前缓存控件下的核心缓存
	nodes      NodePicker //在当前的缓存控件中保存所有的缓存节点
	dataGetter DataGetter //当从缓存中无法获取值的时候，使用调用传入的Getter获取值
}

type DataGetter interface {
	Get(key string) ([]byte, error)
}

type DataGetterFunc func(key string) ([]byte, error)

func (f DataGetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	cacheNameSpaces = make(map[string]*CacheNameSpace)
	mx              sync.Mutex
)

func InitNameSapce(name string, spaceCap int64, dataGeeter DataGetter) *CacheNameSpace {
	if name == "" || dataGeeter == nil {
		panic("name & dataGeeter is required")
	}
	mx.Lock()
	defer mx.Unlock()
	if GetNameSpace(name) == nil {
		namespace := &CacheNameSpace{
			name:       name,
			cache:      InitCache(spaceCap),
			dataGetter: dataGeeter,
		}

		cacheNameSpaces[name] = namespace
	}

	return GetNameSpace(name)
}

func GetNameSpace(name string) *CacheNameSpace {
	return cacheNameSpaces[name]
}

func (cacheSpace *CacheNameSpace) Get(key string) (CacheView, error) {
	if key == "" {
		return CacheView{}, fmt.Errorf("key is required")
	}

	if value, ok := cacheSpace.cache.Get(key); ok {
		return value.(CacheView), nil
	}

	return cacheSpace.getFromLoaclOrRemoteNode(key)
}

func (cacheSpace *CacheNameSpace) Put(key string, value CacheView) error {
	mx.Lock()
	defer mx.Unlock()
	if key == "" {
		return errors.New("key is required")
	}

	cacheSpace.cache.Put(key, value)

	return nil
}

func (cacheSpace *CacheNameSpace) getFromLoaclOrRemoteNode(key string) (value CacheView, err error) {
	if cacheSpace.nodes != nil {
		if node, ok := cacheSpace.nodes.NodeSelect(key); ok {
			if value, err = cacheSpace.getFromRemoteNode(node, key); err == nil {
				return value, nil
			}
		}
	}

	return cacheSpace.getFromLocalNode(key)
}

func (cacheSpace *CacheNameSpace) getFromLocalNode(key string) (CacheView, error) {
	bytes, err := cacheSpace.dataGetter.Get(key)
	if err != nil {
		return CacheView{}, err
	}
	value := CacheView{cloneCache(bytes)}
	cacheSpace.refreshIntoLocalCache(key, value)
	return value, nil
}

func (cacheSpace *CacheNameSpace) refreshIntoLocalCache(key string, cache CacheView) {
	cacheSpace.cache.Put(key, cache)
}

func (cacheSpace *CacheNameSpace) getFromRemoteNode(nodeGet NodeGet, key string) (CacheView, error) {
	bytes, err := nodeGet.Get(cacheSpace.name, key)
	if err != nil {
		return CacheView{}, err
	}
	return CacheView{bytes}, nil
}

func (cacheSpace *CacheNameSpace) RegisterCacheNode(nodes NodePicker) {
	cacheSpace.nodes = nodes
}
