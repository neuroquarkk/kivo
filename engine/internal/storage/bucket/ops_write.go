package bucket

import (
	"bytes"
	"kivo/engine/internal/ttl/heap"
	"time"
)

func (b *Bucket) Set(
	key string, value []byte, ttl time.Duration,
) (bool, int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now().UnixNano()
	ent, exists := b.data[key]
	isNew := !exists

	var delta int64
	if isNew {
		delta = entrySize(key, value)
	} else {
		delta = int64(len(value)) - int64(len(ent.value))
	}

	if entrySize(key, value) > b.targetBytes {
		return false, 0, ErrValueTooBig
	}

	var removed int
	if b.currentSize+delta >= b.thresholdBytes {
		var err error
		removed, err = b.evictKeys(key, delta)
		if err != nil {
			return false, removed, err
		}
	}

	if isNew {
		ent = &entry{}
		b.data[key] = ent
		b.keys = append(b.keys, key)
		ent.feq.Store(5)
	} else {
		newScore := ent.feq.Add(1)
		if newScore < 5 {
			ent.feq.Store(5)
		}
	}

	b.currentSize += delta

	var expiresAt int64
	tracked := ent.item != nil

	if ttl > 0 {
		expiresAt = now + int64(ttl)
		if tracked {
			// expiry -> expiry (update)
			b.ttlHeap.Update(ent.item, expiresAt)
		} else {
			// permanent -> expiry (new)
			ent.item = b.ttlHeap.Push(key, expiresAt)
		}
	} else if tracked {
		// expiry -> permanent
		b.ttlHeap.Remove(ent.item)
		ent.item = nil
	}

	ent.value = bytes.Clone(value)
	ent.expiresAt = expiresAt

	return isNew, removed, nil
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
	b.currentSize -= entrySize(key, ent.value)
	return true
}

func (b *Bucket) Flush() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data = make(map[string]*entry)
	b.ttlHeap = heap.New()
	b.currentSize = 0
}
