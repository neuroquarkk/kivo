package bucket

import "sync"

type entry struct {
	value     []byte
	expiresAt int64
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

func (e entry) isExpired(now int64) bool {
	return e.expiresAt != 0 && now > e.expiresAt
}
