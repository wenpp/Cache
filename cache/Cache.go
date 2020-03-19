package cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	//字典表，用于映射key和链表的element
	dict map[string]*list.Element
	//双向链表,链表头部是最长时间没有使用的entry,尾部是最近使用的
	dll *list.List
	//缓存的最大容量
	maxSize int64
	//缓存的当前容量
	currSize int64
}

type Value interface {
	Len() int
}

type Entry struct {
	key   string
	value Value
}

var (
	mu sync.Mutex
	rw sync.RWMutex
)

func (cache *Cache) Len() int {
	return cache.dll.Len()
}

func InitCache(maxSize int64) *Cache {
	return &Cache{
		maxSize: maxSize,
		dll:     list.New(),
		dict:    make(map[string]*list.Element),
	}
}

func (cache *Cache) Put(key string, value Value) {
	mu.Lock()
	defer mu.Unlock()
	if cache == nil {
		cache = InitCache(cache.maxSize)
	}

	if element, ok := cache.dict[key]; ok {
		cache.dll.MoveToFront(element)
		entry := element.Value.(*Entry)
		entry.value = value
	} else {
		element := cache.dll.PushFront(&Entry{key, value})
		cache.dict[key] = element
		cache.currSize += int64(len(key)) + int64(value.Len())
	}

	//初始化的时候如果为0则不限制大小，永不删除最久没有使用的entry
	if cache.maxSize != 0 && cache.maxSize < cache.currSize {
		cache.DeleteOldest()
	}
}

func (cache *Cache) Get(key string) (value Value, ok bool) {
	rw.RLock()
	defer rw.RUnlock()
	if element, ok := cache.dict[key]; ok {
		cache.dll.MoveToBack(element)
		entry := element.Value.(*Entry)
		return entry.value, true
	}
	return
}

func (cache *Cache) DeleteOldest() {
	element := cache.dll.Back()
	if element != nil {
		cache.dll.Remove(element)
		entry := element.Value.(*Entry)
		delete(cache.dict, entry.key)
		cache.currSize -= int64(len(entry.key)) + int64(entry.value.Len())
	}
}
