package validation

import (
	"errors"
	"strings"
)

var (
	ErrEmptyKey = errors.New("key is empty")
)

func ValidateKey(key string) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if strings.TrimSpace(key) == "" {
		return ErrEmptyKey
	}

	return nil
}
