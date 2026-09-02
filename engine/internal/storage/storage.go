package storage

import (
	"hash/maphash"
	"kivo/engine/internal/storage/bucket"
)

const numBucket = 256

type Store struct {
	buckets []*bucket.Bucket
	seed    maphash.Seed
}

func New() *Store {
	buckets := make([]*bucket.Bucket, numBucket)
	for i := range numBucket {
		buckets[i] = bucket.New()
	}

	return &Store{
		buckets: buckets,
		seed:    maphash.MakeSeed(),
	}
}

func (s *Store) getBucket(key string) *bucket.Bucket {
	h := maphash.String(s.seed, key)
	idx := h & (numBucket - 1)
	return s.buckets[idx]
}
