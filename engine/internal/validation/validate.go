package validation

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyKey      = errors.New("key is empty")
	ErrKeyTooLarge   = errors.New("key is too large")
	ErrValueTooLarge = errors.New("value is too large")
	ErrNegativeTTL   = errors.New("ttl cannot be negative")
)

func CheckKey(key string, max int) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if strings.TrimSpace(key) == "" {
		return ErrEmptyKey
	}

	if max > 0 && len(key) > max {
		return ErrKeyTooLarge
	}

	return nil
}

func CheckValue(value []byte, max int) error {
	if max > 0 && len(value) > max {
		return ErrValueTooLarge
	}

	return nil
}

func CheckTTL(ttl time.Duration) error {
	if ttl < 0 {
		return ErrNegativeTTL
	}
	return nil
}
