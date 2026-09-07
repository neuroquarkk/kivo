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

func (b *Bucket) evictKeys(ignoreKey string) int {
	targetSize := int64(float64(b.maxSize) * evictionTargetRatio)
	now := time.Now().UnixNano()

	var evictedCount int
	for b.currentSize > targetSize {
		var lowestKey string
		var lowestFeq uint32 = math.MaxUint32
		var sampled int

		for k, e := range b.data {
			if k == ignoreKey {
				continue
			}

			if now-e.createdAt < int64(evictionGracePeriod) {
				continue
			}

			feq := e.feq.Load()
			if feq < lowestFeq {
				lowestFeq = feq
				lowestKey = k
			}
			sampled++
			if sampled >= evictionSampleSize {
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
		b.currentSize -= int64(len(lowestKey))
		b.currentSize -= int64(len(ent.value))
		evictedCount++
	}

	return evictedCount
}
