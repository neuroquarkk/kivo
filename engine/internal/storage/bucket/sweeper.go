package bucket

import "math"

const (
	decayFactor = 0.5
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

func (b *Bucket) evictKeys(ignoreKey string, delta, now int64) int {
	graceCutoff := now - int64(b.cfg.GracePeriod)
	maxInspect := b.cfg.SampleSize * 5

	var evictedCount int

	for (b.currentSize + delta) > b.targetBytes {
		var (
			lowestKey string
			lowestFeq uint32 = math.MaxUint32
			sampled   int
			visited   int
		)

		for key, ent := range b.data {
			if key == ignoreKey {
				continue
			}

			visited++
			if visited > maxInspect {
				break
			}

			if ent.createdAt > graceCutoff {
				continue
			}

			if feq := ent.feq.Load(); feq < lowestFeq {
				lowestFeq = feq
				lowestKey = key
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

func (b *Bucket) Decay() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ent := range b.data {
		feq := ent.feq.Load()
		decayed := float64(feq) * decayFactor
		ent.feq.Store(uint32(decayed))
	}
}
