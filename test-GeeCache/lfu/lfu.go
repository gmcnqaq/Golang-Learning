package lfu

type Cache struct {
	capacity   int64
	size       int64
	minFreq    int64
	keyToVal   map[int]int
	keyToFreq  map[int]int
	freqToKeys map[int][]int
}

func (c *Cache) Get(key int) int {
	if value, ok := c.keyToVal[key]; ok {
		return value
	} else {
		return -1
	}
}
