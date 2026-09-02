package storage

func (s *Store) Set(key string, value []byte) {
	b := s.getBucket(key)
	b.Set(key, value)
}

func (s *Store) Delete(key string) {
	b := s.getBucket(key)
	b.Delete(key)
}

func (s *Store) Get(key string) ([]byte, error) {
	b := s.getBucket(key)
	return b.Get(key)
}

func (s *Store) Exists(key string) bool {
	b := s.getBucket(key)
	return b.Exists(key)
}
