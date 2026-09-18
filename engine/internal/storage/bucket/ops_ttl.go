package bucket

import "time"

func (b *Bucket) TTL(key string) (time.Duration, bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ent, ok := b.data[key]
	if !ok {
		return 0, false, ErrNotFound
	}

	now := time.Now().UnixNano()
	if ent.isExpired(now) {
		return 0, false, ErrNotFound
	}

	if ent.item == nil {
		return 0, false, nil // permanent
	}

	remaining := time.Duration(ent.expiresAt - now)
	return remaining, true, nil
}

func (b *Bucket) Expire(key string, ttl time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	ent, ok := b.data[key]
	if !ok {
		return ErrNotFound
	}

	now := time.Now().UnixNano()
	if ent.isExpired(now) {
		return ErrNotFound
	}

	expiresAt := now + int64(ttl)
	if ent.item != nil {
		b.ttlHeap.Update(ent.item, expiresAt)
	} else {
		ent.item = b.ttlHeap.Push(key, expiresAt)
	}
	ent.expiresAt = expiresAt

	return nil
}

func (b *Bucket) Persist(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	ent, ok := b.data[key]
	if !ok {
		return ErrNotFound
	}

	now := time.Now().UnixNano()
	if ent.isExpired(now) {
		return ErrNotFound
	}

	if ent.item != nil {
		b.ttlHeap.Remove(ent.item)
		ent.item = nil
		ent.expiresAt = 0
	}

	return nil
}
