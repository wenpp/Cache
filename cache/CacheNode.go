package cache

import (
	"net/http"
	"strings"
)

type CacheNode struct {
	nodeName string
	basePath string
	hashring *HashRing
}

type NodeGet interface {
	Get(group string, key string) ([]byte, error)
}

type NodePicker interface {
	NodeSelect(key string) (nodeGet NodeGet, ok bool)
}

func InitCacheNode(nodeName string) *CacheNode {
	return &CacheNode{
		nodeName: nodeName,
		basePath: DefaultBasePath,
	}
}

func (cacheNode *CacheNode) HandlerHttpRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path[len(cacheNode.basePath):], "/", 2)

	nameSpaceName := parts[0]

	key := parts[1]

	namespace := GetNameSpace(nameSpaceName)

	if namespace == nil {
		http.Error(w, "No such namespace found"+nameSpaceName, http.StatusInternalServerError)
		return
	}

	cache, err := namespace.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/text")
	w.Write(cache.CacheSlice())
}
