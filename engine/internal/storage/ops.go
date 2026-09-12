package storage

import "time"

func (s *Store) Set(key string, value []byte, ttl time.Duration) {
	b := s.getBucket(key)

	isNew, removed := b.Set(key, value, ttl)
	if isNew {
		s.stats.keyCount.Add(1)
	}
	s.stats.sets.Add(1)
	if removed > 0 {
		s.stats.keyCount.Add(-int64(removed))
		s.stats.evictions.Add(int64(removed))
	}
}

func (s *Store) Delete(key string) {
	b := s.getBucket(key)

	deleted := b.Delete(key)
	if deleted {
		s.stats.keyCount.Add(-1)
	}
	s.stats.deletes.Add(1)
}

func (s *Store) Flush() {
	for _, b := range s.buckets {
		b.Flush()
	}
	s.stats.keyCount.Store(0)
}

func (s *Store) Get(key string) ([]byte, error) {
	b := s.getBucket(key)

	value, err := b.Get(key)
	if err != nil {
		s.stats.misses.Add(1)
		return nil, err
	}

	s.stats.hits.Add(1)
	return value, nil
}

func (s *Store) Exists(key string) bool {
	b := s.getBucket(key)
	exists := b.Exists(key)

	if exists {
		s.stats.hits.Add(1)
	} else {
		s.stats.misses.Add(1)
	}

	return exists
}

func (s *Store) Count() int64 {
	return s.stats.keyCount.Load()
}

func (s *Store) Info() StatsSnapshot {
	return StatsSnapshot{
		KeyCount:     s.stats.keyCount.Load(),
		Sets:         s.stats.sets.Load(),
		Deletes:      s.stats.deletes.Load(),
		Hits:         s.stats.hits.Load(),
		Misses:       s.stats.misses.Load(),
		Evictions:    s.stats.evictions.Load(),
		MemLimit:     s.memLimit,
		MemPerBucket: s.memLimit / int64(len(s.buckets)),
	}
}
