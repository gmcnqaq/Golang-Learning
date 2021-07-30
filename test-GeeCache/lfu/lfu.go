package lfu

type Cache struct {
	capacity   int64
	size       int64
	minFreq    int64
	keyToVal   map[string]int
	keyToFreq  map[string]int
	freqToKeys map[string][]int
}

func (c *Cache) Get(key string) int {
	if value, ok := c.keyToVal[key]; ok {
		return value
	} else {
		return -1
	}
}
