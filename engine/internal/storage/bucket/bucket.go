package bucket

import (
	"sync"
	"sync/atomic"
	"time"

	"kivo/engine/internal/ttl/heap"
)

type Config struct {
	MaxSize        int64
	ThresholdRatio float64
	TargetRatio    float64
	SampleSize     int
	GracePeriod    time.Duration
}

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
	cfg         Config
	currentSize int64
}

func New(cfg Config) *Bucket {
	return &Bucket{
		data:    make(map[string]*entry),
		ttlHeap: heap.New(),
		cfg:     cfg,
	}
}

func (e *entry) isExpired(now int64) bool {
	return e.expiresAt != 0 && now > e.expiresAt
}

func entrySize(key string, val []byte) int64 {
	return int64(len(key) + len(val))
}
