package bucket

func (b *Bucket) Set(key string, value []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data[key] = entry{value: value}
}

func (b *Bucket) Delete(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.data, key)
}
