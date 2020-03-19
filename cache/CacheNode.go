package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
)

type CacheNode struct {
	nodeAddr         string
	basePath         string
	hashring         *HashRing
	cacheCoordinater map[string]*CacheCoordinator
}

type CacheCoordinator struct {
	URL string
}

type NodeGet interface {
	Get(group string, key string) ([]byte, error)
}

type NodePicker interface {
	NodeSelect(key string) (nodeGet NodeGet, ok bool)
}

func (cacheNode *CacheNode) NodeSelect(key string) (NodeGet, bool) {
	if node := cacheNode.hashring.GetHashKey(key); node != "" && node != cacheNode.nodeAddr {
		return cacheNode.cacheCoordinater[node], true
	}
	return nil, false
}

func (node *CacheCoordinator) Get(group string, key string) ([]byte, error) {
	url := node.URL + url2.QueryEscape(group) + url2.QueryEscape(key)

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server return error")
	}

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("Response Body Error: %v", err)
	}

	return bytes, nil

}

func InitCacheNodes(nodeAddr string) *CacheNode {
	return &CacheNode{
		nodeAddr: nodeAddr,
		basePath: DefaultPath,
	}
}

func (cacheNode *CacheNode) UpdateNodePool(nodes ...string) {
	cacheNode.hashring = InitHashRing(defaultReplicasNode)
	cacheNode.hashring.Add(nodes...)
	cacheNode.cacheCoordinater = make(map[string]*CacheCoordinator, len(nodes))

	for _, node := range nodes {
		cacheNode.cacheCoordinater[node] = &CacheCoordinator{URL: node + cacheNode.basePath}
	}
}

func StartCacheNode(localAddress string, nodesAddress []string, cacheSpace *CacheNameSpace) {
	cacheNode := InitCacheNodes(localAddress)
	cacheNode.UpdateNodePool(nodesAddress...)
	cacheSpace.RegisterCacheNode(cacheNode)
	log.Println("CacheNode:" + localAddress + " is running")
	http.ListenAndServe(localAddress, nil)
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
