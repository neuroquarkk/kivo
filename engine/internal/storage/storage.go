package storage

import (
	"context"
	"hash/maphash"
	"sync/atomic"
	"time"

	"kivo/engine/internal/storage/bucket"
)

const (
	numBucket     = 256
	sweepDuration = 1 * time.Minute
)

type Store struct {
	buckets  []*bucket.Bucket
	seed     maphash.Seed
	keyCount atomic.Int64
}

func New(ctx context.Context) *Store {
	buckets := make([]*bucket.Bucket, numBucket)
	for i := range numBucket {
		buckets[i] = bucket.New()
	}

	s := &Store{
		buckets: buckets,
		seed:    maphash.MakeSeed(),
	}

	s.startSweeper(ctx)
	return s
}

func (s *Store) getBucket(key string) *bucket.Bucket {
	h := maphash.String(s.seed, key)
	idx := h & (numBucket - 1)
	return s.buckets[idx]
}

func (s *Store) startSweeper(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(sweepDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now().UnixNano()
				for _, b := range s.buckets {
					removed := b.SweepExpired(now)
					s.keyCount.Add(-removed)
				}
			}
		}
	}()
}
