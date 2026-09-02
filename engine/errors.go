package engine

import (
	"kivo/engine/internal/storage/bucket"
	"kivo/engine/internal/validation"
)

var (
	ErrEmptyKey = validation.ErrEmptyKey
	ErrNotFound = bucket.ErrNotFound
)
