package engine

import (
	"kivo/engine/internal/storage"
	"kivo/engine/internal/validation"
	"time"
)

type Engine struct {
	store *storage.Store
}

func New() *Engine {
	return &Engine{
		store: storage.New(),
	}
}

func (e *Engine) Set(key string, value []byte, ttl time.Duration) error {
	if err := validation.ValidateKey(key); err != nil {
		return err
	}

	e.store.Set(key, value, ttl)
	return nil
}

func (e *Engine) Delete(key string) error {
	if err := validation.ValidateKey(key); err != nil {
		return err
	}

	e.store.Delete(key)
	return nil
}

func (e *Engine) Get(key string) ([]byte, error) {
	if err := validation.ValidateKey(key); err != nil {
		return nil, err
	}

	value, err := e.store.Get(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (e *Engine) Exists(key string) (bool, error) {
	if err := validation.ValidateKey(key); err != nil {
		return false, err
	}

	exists := e.store.Exists(key)
	return exists, nil
}

func (e *Engine) Count() int64 {
	return e.store.Count()
}
