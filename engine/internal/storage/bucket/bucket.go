package bucket

import "sync"

type entry struct {
	value []byte
}

type Bucket struct {
	mu   sync.RWMutex
	data map[string]entry
}

func New() *Bucket {
	return &Bucket{
		data: make(map[string]entry),
	}
}
