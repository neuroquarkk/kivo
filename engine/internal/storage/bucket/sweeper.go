package bucket

func (b *Bucket) SweepExpired(now int64) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	var removed int64
	for {
		item := b.ttlHeap.Peek()
		if item == nil || now < item.ExpiresAt {
			break
		}

		b.ttlHeap.Pop()

		ent := b.data[item.Key]
		b.currentSize -= entrySize(item.Key, ent.value)
		b.removeKeysAt(ent.idx)
		delete(b.data, item.Key)
		removed++
	}

	return removed
}

func (b *Bucket) evictKeys(ignoreKey string, delta int64) (int, error) {
	var evictedCount int
	startCursor := -1
	progressed := false

	for (b.currentSize + delta) > b.targetBytes {
		if len(b.keys) == 0 {
			return evictedCount, ErrValueTooBig
		}

		b.cursor = b.cursor % len(b.keys)
		// completed a lap with nothing evicted and no score decayed either
		if b.cursor == startCursor && !progressed {
			return evictedCount, ErrValueTooBig
		}
		if startCursor == -1 || b.cursor == startCursor {
			startCursor = b.cursor
			progressed = false
		}

		key := b.keys[b.cursor]
		ent := b.data[key]

		if key == ignoreKey {
			b.cursor++
			continue
		}

		score := ent.feq.Load()
		if score > 0 {
			ent.feq.Store(score / 2)
			progressed = true
			b.cursor++
			continue
		}

		// score == 0: evict
		// removeKeysAt swaps the last key into this slot
		// so the cursor is left pointing at the next key to
		// consider without incrementing
		b.currentSize -= entrySize(key, ent.value)
		delete(b.data, key)
		if ent.item != nil {
			b.ttlHeap.Remove(ent.item)
		}

		b.removeKeysAt(b.cursor)
		evictedCount++
		progressed = true
		startCursor = -1
	}

	return evictedCount, nil
}
