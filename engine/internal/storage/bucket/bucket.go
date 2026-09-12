package bucket

import (
	"errors"
	"sync"
	"sync/atomic"

	"kivo/engine/internal/ttl/heap"
)

var (
	ErrNotFound    = errors.New("key not found")
	ErrValueTooBig = errors.New("value too big for bucket")
)

type Config struct {
	ThresholdBytes int64
	TargetBytes    int64
	GracePeriod    int64
}

type entry struct {
	value     []byte
	expiresAt int64
	item      *heap.Item
	feq       atomic.Uint32
}

type Bucket struct {
	mu      sync.RWMutex
	data    map[string]*entry
	keys    []string
	cursor  int
	ttlHeap *heap.Heap

	currentSize    int64
	thresholdBytes int64
	targetBytes    int64
	gracePeriod    int64
}

func New(cfg Config) *Bucket {
	return &Bucket{
		data:           make(map[string]*entry),
		ttlHeap:        heap.New(),
		thresholdBytes: cfg.ThresholdBytes,
		targetBytes:    cfg.TargetBytes,
		gracePeriod:    cfg.GracePeriod,
	}
}

func (e *entry) isExpired(now int64) bool {
	return e.expiresAt != 0 && now > e.expiresAt
}

func entrySize(key string, val []byte) int64 {
	return int64(len(key) + len(val))
}
