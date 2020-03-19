package test

import (
	"Cache/cache"
	"Cache/db"
	"fmt"
	"reflect"
	"testing"
)

func TestGetGroup(t *testing.T) {
	cache.InitNameSapce("user", 0, cache.DataGetterFunc(
		func(key string) (bytes []byte, err error) { return }))

	if cache.GetNameSpace("car") != nil {
		t.Fatal("expect nil but got car")
	}

	if cache.GetNameSpace("user") == nil {
		t.Fatal("expect user existd but got nil")
	}
}

func TestInitNameSpaceAndGet(t *testing.T) {
	nsp := cache.InitNameSapce("user", 0, cache.DataGetterFunc(
		func(key string) ([]byte, error) {
			if value, ok := db.DB[key]; ok {
				return []byte(value), nil
			}
			return nil, fmt.Errorf("key not exists")
		}))
	err := nsp.Put("key1", cache.CacheView{[]byte("hello")})
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := cache.CacheView{[]byte("hello")}

	if value, _ := nsp.Get("key1"); !reflect.DeepEqual(value, expect) {
		t.Fatal("expect value is hello actually got", value)
	}
}

func TestGetFromDB(t *testing.T) {
	nsp := cache.InitNameSapce("user", 0, cache.DataGetterFunc(
		func(key string) ([]byte, error) {
			if value, ok := db.DB[key]; ok {
				return []byte(value), nil
			}
			return nil, fmt.Errorf("key not exists")
		}))

	expect := cache.CacheView{[]byte("25")}

	//BMW is 75
	if value, _ := nsp.Get("BMW"); !reflect.DeepEqual(value, expect) {
		t.Fatalf("expect value of key %s is %s but actually got %s", "BMW", expect, value)
	}
}
