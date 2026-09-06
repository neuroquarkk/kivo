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

type StatsSnapshot struct {
	KeyCount     int64
	Sets         int64
	Deletes      int64
	Hits         int64
	Misses       int64
	Evictions    int64
	MemLimit     int64
	MemPerBucket int64
}

type stats struct {
	keyCount  atomic.Int64
	sets      atomic.Int64
	deletes   atomic.Int64
	hits      atomic.Int64
	misses    atomic.Int64
	evictions atomic.Int64
}

type Store struct {
	buckets  []*bucket.Bucket
	seed     maphash.Seed
	stats    stats
	memLimit int64
}

func New(ctx context.Context, memLimit int64) *Store {
	perBucketLimit := memLimit / int64(numBucket)
	buckets := make([]*bucket.Bucket, numBucket)
	for i := range numBucket {
		buckets[i] = bucket.New(perBucketLimit)
	}

	s := &Store{
		buckets:  buckets,
		seed:     maphash.MakeSeed(),
		memLimit: memLimit,
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
					s.stats.keyCount.Add(-removed)
					s.stats.evictions.Add(removed)
				}
			}
		}
	}()
}
