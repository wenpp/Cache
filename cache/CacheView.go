package cache

type CacheView struct {
	CacheByte []byte
}

func (c CacheView) CacheSlice() []byte {
	return cloneCache(c.CacheByte)
}

func cloneCache(cacheByte []byte) []byte {
	c := make([]byte, len(cacheByte))
	copy(c, cacheByte)
	return c
}

func (c CacheView) String() string {
	return string(c.CacheByte)
}

func (c CacheView) Len() int {
	return len(c.CacheByte)
}
