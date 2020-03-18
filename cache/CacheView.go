package cache

type CacheView struct {
	cacheByte []byte
}

func (c CacheView) CacheSlice() []byte {
	return cloneCache(c.cacheByte)
}

func cloneCache(cacheByte []byte) []byte {
	c := make([]byte, len(cacheByte))
	copy(c, cacheByte)
	return c
}
