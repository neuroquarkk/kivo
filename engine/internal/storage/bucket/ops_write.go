package bucket

import (
	"bytes"
	"time"
)

func (b *Bucket) Set(key string, value []byte, ttl time.Duration) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	ent, ok := b.data[key]

	hasTTL := ttl > 0
	tracked := ent.item != nil

	var expiresAt int64
	if hasTTL {
		expiresAt = time.Now().Add(ttl).UnixNano()
	}

	switch {
	case hasTTL && !tracked:
		// permanent -> expiry (new)
		ent.item = b.ttlHeap.Push(key, expiresAt)
	case hasTTL && tracked:
		// expiry -> expiry (update)
		b.ttlHeap.Update(ent.item, expiresAt)
	case !hasTTL && tracked:
		// expiry -> permanent
		b.ttlHeap.Remove(ent.item)
		ent.item = nil
	}

	ent.value = bytes.Clone(value)
	ent.expiresAt = expiresAt
	b.data[key] = ent

	kLen := int64(len(key))
	vLen := int64(len(value))

	if !ok {
		b.currentSize += kLen + vLen
	} else {
		oldVLen := int64(len(ent.value))
		b.currentSize -= oldVLen
		b.currentSize += vLen
	}
	return !ok
}

func (b *Bucket) Delete(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	ent, ok := b.data[key]
	if !ok {
		return false
	}

	if ent.item != nil {
		b.ttlHeap.Remove(ent.item)
	}

	delete(b.data, key)

	b.currentSize -= int64(len(key))
	b.currentSize -= int64(len(ent.value))
	return true
}
