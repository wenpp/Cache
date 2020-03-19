package test

import (
	cache2 "Cache/cache"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestPut(t *testing.T) {
	cache := cache2.InitCache(int64(0))
	cache.Put("key", String("Hello"))

	expect := 1

	if expect != cache.Len() {
		t.Fatal("expected 1 but got", cache.Len())
	}
}

func TestGet(t *testing.T) {
	cache := cache2.InitCache(int64(0))
	cache.Put("key1", String("Hello"))
	cache.Put("key2", String("morning"))

	expect1 := String("Hello")
	expect2 := String("morning")

	if entry, _ := cache.Get("key1"); entry != expect1 {
		t.Fatal("expected Hello but got", entry)
	}

	if entry, _ := cache.Get("key2"); entry != expect2 {
		t.Fatal("expected morning but got", entry)
	}
}

func TestDeleteOldest(t *testing.T) {
	key1 := "key1"
	key2 := "key2"
	key3 := "key3"

	value1 := String("value1")
	value2 := String("value2")
	value3 := String("value=3")

	maxSizeCorrect := int64(len(key1) + len(key2) + value1.Len() + value2.Len())
	cache := cache2.InitCache(maxSizeCorrect)

	cache.Put(key1, value1)
	cache.Put(key2, value2)
	cache.Put(key3, value3)

	if _, ok := cache.Get(key1); ok {
		t.Fatal("Key1 still existed")
	}

	expectEntryLen := 2

	if cache.Len() != expectEntryLen {
		t.Fatal("expect cache size is 2 but got", cache.Len())
	}

}
