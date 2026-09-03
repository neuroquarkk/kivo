package bucket

import (
	"bytes"
	"time"
)

func (b *Bucket) Set(key string, value []byte, ttl time.Duration) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.data[key]

	var expiresAt int64
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl).UnixNano()
	}

	b.data[key] = entry{
		value:     bytes.Clone(value),
		expiresAt: expiresAt,
	}

	return !ok
}

func (b *Bucket) Delete(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.data[key]; !ok {
		return false
	}

	delete(b.data, key)
	return true
}
