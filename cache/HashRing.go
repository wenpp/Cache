package cache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//一致性hash环
type HashRing struct {
	keys           []int          //哈希环
	virtualNodeMap map[int]string //虚拟节点与真实节点映射
	replicas       int            //虚拟节点数量
}

func Init(replicas int) *HashRing {
	ch := &HashRing{
		virtualNodeMap: make(map[int]string),
		replicas:       replicas,
	}
	return ch
}

func (chash *HashRing) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < chash.replicas; i++ {
			hashValue := int(crc32.ChecksumIEEE([]byte(strconv.Itoa(i) + key)))
			chash.keys = append(chash.keys, hashValue)
			chash.virtualNodeMap[hashValue] = key
		}
	}

	sort.Ints(chash.keys)
}

func (chash *HashRing) Get(key string) string {
	if len(chash.keys) == 0 {
		return ""
	}

	hashValue := int(crc32.ChecksumIEEE([]byte(key)))
	index := sort.Search(len(chash.keys), func(i int) bool {
		return chash.keys[i] >= hashValue
	})

	return chash.virtualNodeMap[chash.keys[index%len(chash.keys)]]
}
