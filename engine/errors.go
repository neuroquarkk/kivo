package engine

import (
	"kivo/engine/internal/storage/bucket"
	"kivo/engine/internal/validation"
)

var (
	ErrEmptyKey      = validation.ErrEmptyKey
	ErrKeyTooLarge   = validation.ErrKeyTooLarge
	ErrValueTooLarge = validation.ErrValueTooLarge
	ErrNegativeTTL   = validation.ErrNegativeTTL
	ErrNotFound      = bucket.ErrNotFound
	ErrValueTooBig   = bucket.ErrValueTooBig
)
