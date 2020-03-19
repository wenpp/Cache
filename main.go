package main

import (
	"Cache/cache"
	"Cache/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cacheQueryServer := "127.0.0.1:7000"
	cacheSpace := intiCacheNameSpace()

	nodeList := []string{"localhost:8001", "localhost:8002", "localhost:8003"}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("hello"))
	})

	http.Handle("/cache", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			value, err := cacheSpace.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(value.CacheSlice())
		}))

	LaunchCacheNode(nodeList, cacheSpace)

	log.Println("Query Server is launched successul! At address:", cacheQueryServer)
	http.ListenAndServe(cacheQueryServer, nil)

}

func intiCacheNameSpace() *cache.CacheNameSpace {
	return cache.InitNameSapce("user", 0, cache.DataGetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := db.DB[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("#{key} is not exists")
		}))
}

func LaunchCacheNode(nodeList []string, cacheSpace *cache.CacheNameSpace) {
	for _, node := range nodeList {
		go cache.StartCacheNode(node, nodeList, cacheSpace)
	}
}
