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

	// MaxValueSize is the maximum allowed key size in bytes.
	// Zero means no limit
	MaxValueSize int

	// DefaultTTL is the default TTL for Set calls without a specified TTL
	// Zero means keys are permanent (default)
	DefaultTTL time.Duration

	MemoryLimitMB int

	limitBytes int64
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
	if o.MemoryLimitMB < 0 {
		return errors.New("memory limit cannot be negative")
	}

	if o.MemoryLimitMB == 0 {
		memLimit, err := getTotalMemoryBytes()
		if err != nil {
			return fmt.Errorf("auto-memory discovery failed: %w", err)
		}
		o.limitBytes = memLimit
	} else {
		o.limitBytes = int64(o.MemoryLimitMB) * 1024 * 1024
	}

	return nil
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
