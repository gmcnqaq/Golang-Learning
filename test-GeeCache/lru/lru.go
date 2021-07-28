package lru

type Cache struct {
	cache    map[string]interface{}
	capacity int64
	size     int64
}
