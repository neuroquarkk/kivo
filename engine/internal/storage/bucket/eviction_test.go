package bucket

import (
	"errors"
	"testing"
	"time"
)

func TestEviction(t *testing.T) {
	// normal eviction, fill past threshold
	b := New(Config{ThresholdBytes: 500, TargetBytes: 400})
	val := make([]byte, 150)

	b.Set("k1", val, 0)
	b.Set("k2", val, 0)
	b.Set("k3", val, 0)
	_, removed, err := b.Set("k4", val, 0)
	if err != nil || removed == 0 {
		t.Fatalf("expected an eviction, got removed=%d err=%v", removed, err)
	}

	// value too big outright
	_, _, err = b.Set("big", make([]byte, 1000), 0)
	if !errors.Is(err, ErrValueTooBig) {
		t.Fatalf("expected ErrValueTooBig, got %v, key %s", err, "big")
	}

	// only key present, unevictable
	b2 := New(Config{ThresholdBytes: 500, TargetBytes: 400})
	b2.data["solo"] = &entry{value: make([]byte, 300)}
	b2.keys = []string{"solo"}
	b2.currentSize = 304

	// must error, not hang
	done := make(chan error, 1)
	go func() {
		_, err := b2.evictKeys("solo", 200) // pretend an incoming big write
		done <- err
	}()

	select {
	case err := <-done:
		if !errors.Is(err, ErrValueTooBig) {
			t.Fatalf("expected ErrValueTooBig, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("evictKeys hung")
	}
}
