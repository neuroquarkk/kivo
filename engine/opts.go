package engine

import (
	"errors"
	"time"
)

type Opts struct {
	// MaxKeySize is the maximum allowed key size in bytes.
	// Zero means no limit
	MaxKeySize int

	// MaxValueSize is the maximum allowed key size in bytes.
	// Zero means no limit
	MaxValueSize int

	// DefaultTTL is the default TTL for Set calls without a specified TTL
	// Zero means keys are permanent (default)
	DefaultTTL time.Duration
}

func DefaultOpts() *Opts {
	return &Opts{}
}

func (o *Opts) prepare() error {
	// validate
	if o.MaxKeySize < 0 {
		return errors.New("key size cannot be negative")
	}
	if o.MaxValueSize < 0 {
		return errors.New("value size cannot be negative")
	}
	if o.DefaultTTL < 0 {
		return errors.New("default TTL cannot be negative")
	}

	return nil
}
