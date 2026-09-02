package bucket

import "bytes"

func (b *Bucket) Set(key string, value []byte) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.data[key]

	b.data[key] = entry{
		value: bytes.Clone(value),
	}

	return !ok
}

func (b *Bucket) Delete(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.data[key]; !ok {
		return false
	}

	delete(b.data, key)
	return true
}
