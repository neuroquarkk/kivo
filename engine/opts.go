package engine

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type Opts struct {
	// MaxKeySize is the maximum allowed key size in bytes.
	// Zero means no limit
	MaxKeySize int

	// MaxValueSize is the maximum allowed value size in bytes
	// Zero means no limit
	MaxValueSize int

	// DefaultTTL is the default TTL for Set calls without a specified TTL
	// Zero means keys are permanent (default)
	DefaultTTL time.Duration

	// MemoryLimitMB is the maximum memory limit in megabytes for the store
	// Zero defaults to 80% of available host RAM
	MemoryLimitMB int

	// EvictionThresholdRatio is the capacity ratio (0.0 to 1.0) at which bucket evictions trigger
	// Zero defaults to 0.99 (99%)
	EvictionThresholdRatio float64

	// EvictionTargetRatio is the target capacity ratio (0.0 to 1.0) to which buckets clear down during eviction
	// Zero defaults to 0.95 (95%)
	EvictionTargetRatio float64

	// EvictionSampleSize is the number of random keys sampled during an eviction pass
	// Zero defaults to 5
	EvictionSampleSize int

	// EvictionGracePeriod is the minimum age required for a key to be eligible for eviction sampling
	// Zero defaults to 1.5 seconds
	EvictionGracePeriod time.Duration
}

type internalOpts struct {
	limitBytes     int64
	thresholdRatio float64
	targetRatio    float64
	sampleSize     int
	gracePeriod    time.Duration
	maxKeySize     int
	maxValueSize   int
	defaultTTL     time.Duration
}

func DefaultOpts() *Opts {
	return &Opts{
		MaxKeySize:             256,
		MaxValueSize:           512 * 1024,
		MemoryLimitMB:          0,
		EvictionThresholdRatio: 0.99,
		EvictionTargetRatio:    0.95,
		EvictionSampleSize:     5,
		EvictionGracePeriod:    1500 * time.Millisecond,
	}
}

func (o *Opts) prepare() (internalOpts, error) {
	if o.MaxKeySize < 0 {
		return internalOpts{}, errors.New("key size cannot be negative")
	}
	if o.MaxValueSize < 0 {
		return internalOpts{}, errors.New("value size cannot be negative")
	}
	if o.DefaultTTL < 0 {
		return internalOpts{}, errors.New("default TTL cannot be negative")
	}
	if o.MemoryLimitMB < 0 {
		return internalOpts{}, errors.New("memory limit cannot be negative")
	}
	if o.EvictionThresholdRatio < 0 || o.EvictionThresholdRatio > 1.0 {
		return internalOpts{}, errors.New("eviction threshold ratio must be between 0 and 1.0")
	}
	if o.EvictionTargetRatio < 0 || (o.EvictionThresholdRatio > 0 && o.EvictionTargetRatio >= o.EvictionThresholdRatio) {
		return internalOpts{}, errors.New("eviction target ratio must be less than threshold ratio")
	}
	if o.EvictionSampleSize < 0 {
		return internalOpts{}, errors.New("eviction sample size cannot be negative")
	}
	if o.EvictionGracePeriod < 0 {
		return internalOpts{}, errors.New("eviction grace period cannot be negative")
	}

	var limitBytes int64
	if o.MemoryLimitMB == 0 {
		memLimit, err := getTotalMemoryBytes()
		if err != nil {
			return internalOpts{}, fmt.Errorf("auto-memory discovery failed: %w", err)
		}
		limitBytes = memLimit
	} else {
		limitBytes = int64(o.MemoryLimitMB) * 1024 * 1024
	}

	threshold := o.EvictionThresholdRatio
	if threshold == 0 {
		threshold = 0.99
	}

	target := o.EvictionTargetRatio
	if target == 0 {
		target = 0.95
	}

	sampleSize := o.EvictionSampleSize
	if sampleSize == 0 {
		sampleSize = 5
	}

	gracePeriod := o.EvictionGracePeriod
	if gracePeriod == 0 {
		gracePeriod = 1500 * time.Millisecond
	}

	return internalOpts{
		limitBytes:     limitBytes,
		thresholdRatio: threshold,
		targetRatio:    target,
		sampleSize:     sampleSize,
		gracePeriod:    gracePeriod,
		maxKeySize:     o.MaxKeySize,
		maxValueSize:   o.MaxValueSize,
		defaultTTL:     o.DefaultTTL,
	}, nil
}

func getTotalMemoryBytes() (int64, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}

		var memKB int64
		if _, err := fmt.Sscanf(line, "MemTotal: %d kB", &memKB); err == nil {
			totalBytes := memKB * 1024
			safeLimit := float64(totalBytes) * 0.80
			return int64(safeLimit), nil
		}
	}

	return 0, errors.New("MemTotal not found in /proc/meminfo")
}
