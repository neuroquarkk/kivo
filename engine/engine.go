package engine

import (
	"context"
	"time"

	"kivo/engine/internal/storage"
	"kivo/engine/internal/validation"
)

type Engine struct {
	store *storage.Store
	opts  *Opts
}

func New(ctx context.Context, opts *Opts) (*Engine, error) {
	if opts == nil {
		opts = &Opts{}
	}

	if err := opts.prepare(); err != nil {
		return nil, err
	}

	e := &Engine{}
	e.store = storage.New(ctx)
	e.opts = opts

	return e, nil
}

func (e *Engine) Set(key string, value []byte, ttl time.Duration) error {
	if err := validation.CheckKey(key, e.opts.MaxKeySize); err != nil {
		return err
	}
	if err := validation.CheckValue(value, e.opts.MaxValueSize); err != nil {
		return err
	}

	if ttl == 0 {
		ttl = e.opts.DefaultTTL
	}

	e.store.Set(key, value, ttl)
	return nil
}

func (e *Engine) Delete(key string) error {
	if err := validation.CheckKey(key, e.opts.MaxKeySize); err != nil {
		return err
	}

	e.store.Delete(key)
	return nil
}

func (e *Engine) Get(key string) ([]byte, error) {
	if err := validation.CheckKey(key, e.opts.MaxKeySize); err != nil {
		return nil, err
	}

	value, err := e.store.Get(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (e *Engine) Exists(key string) (bool, error) {
	if err := validation.CheckKey(key, e.opts.MaxKeySize); err != nil {
		return false, err
	}

	exists := e.store.Exists(key)
	return exists, nil
}

func (e *Engine) Count() int64 {
	return e.store.Count()
}

func (e *Engine) Info() storage.StatsSnapshot {
	return e.store.Info()
}
