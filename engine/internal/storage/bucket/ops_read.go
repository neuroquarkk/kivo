package bucket

import (
	"bytes"
	"errors"
)

var (
	ErrNotFound = errors.New("key not found")
)

func (b *Bucket) Get(key string) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ent, ok := b.data[key]
	if !ok {
		return nil, ErrNotFound
	}

	return bytes.Clone(ent.value), nil
}

func (b *Bucket) Exists(key string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, ok := b.data[key]
	return ok
}
