package bucket

import "math"

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

func (b *Bucket) evictKeys(ignoreKey string) {
	targetSize := int64(float64(b.maxSize) * evictionTargetRatio)

	for b.currentSize > targetSize {
		var lowestKey string
		var lowestFeq uint32 = math.MaxUint32
		var sampled int

		for k, e := range b.data {
			if k == ignoreKey {
				continue
			}

			feq := e.feq.Load()
			if feq < lowestFeq {
				lowestFeq = feq
				lowestKey = k
			}
			sampled++
			if sampled >= 5 {
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
	}
}
