package bucket

import (
	"sync"

	"kivo/engine/internal/ttl/heap"
)

type entry struct {
	value     []byte
	expiresAt int64
	item      *heap.Item
}

type Bucket struct {
	mu      sync.RWMutex
	data    map[string]entry
	ttlHeap *heap.Heap
}

func New() *Bucket {
	return &Bucket{
		data:    make(map[string]entry),
		ttlHeap: heap.New(),
	}
}

func (e entry) isExpired(now int64) bool {
	return e.expiresAt != 0 && now > e.expiresAt
}
