package bucket

import (
	"sync"
	"sync/atomic"
	"time"

	"kivo/engine/internal/ttl/heap"
)

const (
	evictionThresholdRatio = 0.99
	evictionTargetRatio    = 0.95
	evictionSampleSize     = 5
	evictionGracePeriod    = 1500 * time.Millisecond
)

type entry struct {
	value     []byte
	expiresAt int64
	item      *heap.Item
	feq       atomic.Uint32
	createdAt int64
}

type Bucket struct {
	mu          sync.RWMutex
	data        map[string]*entry
	ttlHeap     *heap.Heap
	maxSize     int64
	currentSize int64
}

func New(maxSize int64) *Bucket {
	return &Bucket{
		data:    make(map[string]*entry),
		ttlHeap: heap.New(),
		maxSize: maxSize,
	}
}

func (e *entry) isExpired(now int64) bool {
	return e.expiresAt != 0 && now > e.expiresAt
}
