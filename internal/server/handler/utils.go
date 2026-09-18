package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"kivo/engine"
)

func sendResponse(w http.ResponseWriter, code int, data any) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func sendError(w http.ResponseWriter, code int, msg string) {
	sendResponse(w, code, map[string]string{
		"error": msg,
	})
}

func validateKey(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if len(key) > engine.DefaultMaxKeySize {
		return fmt.Errorf("key exceeds maximum length of %d bytes",
			engine.DefaultMaxKeySize)
	}
	return nil
}

func validateValue(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("value cannot be empty")
	}

	if len(value) > engine.DefaultMaxValueSize {
		return fmt.Errorf("value exceeds maximum length of %d bytes",
			engine.DefaultMaxValueSize)
	}
	return nil
}

func validateTTL(ttl time.Duration) error {
	if ttl < 0 {
		return fmt.Errorf("TTL cannot be negative")
	}
	return nil
}

func parseEngineError(err error) (int, string) {
	switch err {
	case engine.ErrNotFound:
		return http.StatusNotFound, "key not found"
	case engine.ErrEmptyKey:
		return http.StatusBadRequest, "key cannot be empty"
	case engine.ErrKeyTooLarge:
		return http.StatusBadRequest, "key exceeds maximum bytes"
	case engine.ErrValueTooLarge, engine.ErrValueTooBig:
		return http.StatusBadRequest, "value exceeds maximum bytes"
	case engine.ErrNegativeTTL:
		return http.StatusBadRequest, "TTL cannot be negative"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func parseTTL(ttl string) (time.Duration, error) {
	if ttl == "" {
		return 0, nil
	}
	parsed, err := time.ParseDuration(ttl)
	if err != nil {
		return 0, fmt.Errorf("ttl must be a duration string")
	}
	return parsed, nil
}
