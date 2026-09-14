package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"kivo/engine"
)

func sendResponse(w http.ResponseWriter, code int, data map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
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
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}
	if len(value) > engine.DefaultMaxValueSize {
		return fmt.Errorf("value exceeds maximum length of %d bytes",
			engine.DefaultMaxValueSize)
	}
	return nil
}
