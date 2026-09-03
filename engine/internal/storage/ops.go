package storage

import "time"

func (s *Store) Set(key string, value []byte, ttl time.Duration) {
	b := s.getBucket(key)

	isNew := b.Set(key, value, ttl)
	if isNew {
		s.keyCount.Add(1)
	}
}

func (s *Store) Delete(key string) {
	b := s.getBucket(key)

	deleted := b.Delete(key)
	if deleted {
		s.keyCount.Add(-1)
	}
}

func (s *Store) Get(key string) ([]byte, error) {
	b := s.getBucket(key)
	return b.Get(key)
}

func (s *Store) Exists(key string) bool {
	b := s.getBucket(key)
	return b.Exists(key)
}

func (s *Store) Count() int64 {
	return s.keyCount.Load()
}
