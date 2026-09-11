package engine

import (
	"context"
	"time"

	"kivo/engine/internal/storage"
	"kivo/engine/internal/validation"
)

type Engine struct {
	store *storage.Store
	opts  internalOpts
}

func New(ctx context.Context, opts *Opts) (*Engine, error) {
	iOpts, err := opts.prepare()
	if err != nil {
		return nil, err
	}

	e := &Engine{}
	e.store = storage.New(ctx, iOpts.limitBytes)
	e.opts = iOpts

	return e, nil
}

func (e *Engine) Set(key string, value []byte, ttl time.Duration) error {
	if err := validation.CheckKey(key, DefaultMaxKeySize); err != nil {
		return err
	}
	if err := validation.CheckValue(value, DefaultMaxValueSize); err != nil {
		return err
	}
	if err := validation.CheckTTL(ttl); err != nil {
		return err
	}

	if ttl == 0 {
		ttl = e.opts.defaultTTL
	}

	e.store.Set(key, value, ttl)
	return nil
}

func (e *Engine) Delete(key string) error {
	if err := validation.CheckKey(key, DefaultMaxKeySize); err != nil {
		return err
	}

	e.store.Delete(key)
	return nil
}

func (e *Engine) Get(key string) ([]byte, error) {
	if err := validation.CheckKey(key, DefaultMaxKeySize); err != nil {
		return nil, err
	}

	value, err := e.store.Get(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (e *Engine) Exists(key string) (bool, error) {
	if err := validation.CheckKey(key, DefaultMaxKeySize); err != nil {
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
