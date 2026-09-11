package engine

import (
	"fmt"
	"time"

	"kivo/engine/internal/sysmem"
	"kivo/engine/internal/validation"
)

const (
	DefaultMaxKeySize   = 256
	DefaultMaxValueSize = 1 * 1024 * 1024 // 1 MB
)

type Opts struct {
	// MemoryLimitMB is the maximum memory limit in megabytes for the store.
	// Zero defaults to 60% of available host RAM
	MemoryLimitMB uint

	// DefaultTTL is the default TTL for Set calls without a specified TTL.
	// Zero means keys are permanent (default)
	DefaultTTL time.Duration
}

type internalOpts struct {
	limitBytes int64
	defaultTTL time.Duration
}

func DefaultOpts() *Opts {
	return &Opts{
		MemoryLimitMB: 0,
	}
}

func (o *Opts) prepare() (internalOpts, error) {
	if o == nil {
		o = DefaultOpts()
	}
	if err := validation.CheckTTL(o.DefaultTTL); err != nil {
		return internalOpts{}, fmt.Errorf("default TTL invalid: %w", err)
	}

	var limitBytes int64
	if o.MemoryLimitMB == 0 {
		memLimit, err := sysmem.TotalMemory()
		if err != nil {
			return internalOpts{}, fmt.Errorf("auto-memory discovery failed: %w", err)
		}
		limitBytes = memLimit
	} else {
		limitBytes = int64(o.MemoryLimitMB) * 1024 * 1024
	}

	return internalOpts{
		limitBytes: limitBytes,
		defaultTTL: o.DefaultTTL,
	}, nil
}
