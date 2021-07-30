package consistehash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := NewMap(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	hash.Add("6", "4", "2")
}
