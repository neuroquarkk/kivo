package bucket

import (
	"math"
	"time"
)

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
		delete(b.data, item.Key)
		removed++
	}

	return removed
}

func (b *Bucket) evictKeys(ignoreKey string, delta int64) int {
	now := time.Now().UnixNano()

	var evictedCount int
	for (b.currentSize + delta) > b.targetBytes {
		var lowestKey string
		var lowestFeq uint32 = math.MaxUint32
		var sampled int

		for k, e := range b.data {
			if k == ignoreKey {
				continue
			}

			if now-e.createdAt < int64(b.cfg.GracePeriod) {
				continue
			}

			feq := e.feq.Load()
			if feq < lowestFeq {
				lowestFeq = feq
				lowestKey = k
			}
			sampled++
			if sampled >= b.cfg.SampleSize {
				break
			}
		}

		if lowestKey == "" {
			break
		}

		ent := b.data[lowestKey]
		if ent.item != nil {
			b.ttlHeap.Remove(ent.item)
		}

		delete(b.data, lowestKey)
		b.currentSize -= entrySize(lowestKey, ent.value)
		evictedCount++
	}

	return evictedCount
}
