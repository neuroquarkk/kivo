package bucket

import (
	"bytes"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("key not found")
)

func (b *Bucket) Get(key string) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ent, ok := b.data[key]
	if !ok {
		return nil, ErrNotFound
	}

	now := time.Now().UnixNano()
	if ent.isExpired(now) {
		return nil, ErrNotFound
	}

	return bytes.Clone(ent.value), nil
}

func (b *Bucket) Exists(key string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ent, ok := b.data[key]
	if !ok {
		return false
	}

	now := time.Now().UnixNano()
	if ent.isExpired(now) {
		return false
	}

	return true
}
