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
	decayDuration = 10 * time.Minute
)

type Config struct {
	MemLimit       int64
	ThresholdRatio float64
	TargetRatio    float64
	SampleSize     int
	GracePeriod    time.Duration
}

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

func New(ctx context.Context, cfg Config) *Store {
	perBucketLimit := cfg.MemLimit / int64(numBucket)

	bCfg := bucket.Config{
		ThresholdBytes: int64(float64(perBucketLimit) * cfg.ThresholdRatio),
		TargetBytes:    int64(float64(perBucketLimit) * cfg.TargetRatio),
		SampleSize:     cfg.SampleSize,
		GracePeriod:    int64(cfg.GracePeriod),
	}

	buckets := make([]*bucket.Bucket, numBucket)
	for i := range numBucket {
		buckets[i] = bucket.New(bCfg)
	}

	s := &Store{
		buckets:  buckets,
		seed:     maphash.MakeSeed(),
		memLimit: cfg.MemLimit,
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
		sweepTicker := time.NewTicker(sweepDuration)
		decayTicker := time.NewTicker(decayDuration)
		defer sweepTicker.Stop()
		defer decayTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-sweepTicker.C:
				now := time.Now().UnixNano()
				var removed int64
				for _, b := range s.buckets {
					removed += b.SweepExpired(now)
				}
				if removed > 0 {
					s.stats.keyCount.Add(-removed)
					s.stats.evictions.Add(removed)
				}
			case <-decayTicker.C:
				for _, b := range s.buckets {
					b.Decay()
				}
			}
		}
	}()
}
