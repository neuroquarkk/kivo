package bucket

func (b *Bucket) SweepExpired(now int64) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	var removed int64
	for {
		item := b.ttlHeap.Peek()
		if item == nil {
			break
		}

		if now < item.ExpiresAt {
			break
		}

		b.ttlHeap.Pop()
		delete(b.data, item.Key)
		removed++
	}

	return removed
}
