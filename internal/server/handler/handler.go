package handler

import (
	"encoding/json"
	"net/http"

	"kivo/engine"
)

type Handler struct {
	eng *engine.Engine
}

func New(eng *engine.Engine) *Handler {
	return &Handler{eng}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var body PutReq
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := validateValue(body.Value); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ttl, err := parseTTL(body.TTL)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateTTL(ttl); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.eng.Set(key, []byte(body.Value), ttl); err != nil {
		code, msg := parseEngineError(err)
		sendError(w, code, msg)
		return
	}

	sendResponse(w, http.StatusCreated, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.eng.Delete(key); err != nil {
		code, msg := parseEngineError(err)
		sendError(w, code, msg)
		return
	}

	sendResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	value, err := h.eng.Get(key)
	if err != nil {
		code, msg := parseEngineError(err)
		sendError(w, code, msg)
		return
	}

	sendResponse(w, http.StatusOK, map[string]string{
		"key":   key,
		"value": string(value),
	})
}

func (h *Handler) Exists(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	exists, err := h.eng.Exists(key)
	if err != nil {
		code, msg := parseEngineError(err)
		sendError(w, code, msg)
		return
	}

	sendResponse(w, http.StatusOK, map[string]any{
		"key":    key,
		"exists": exists,
	})
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	info := h.eng.Info()
	sendResponse(w, http.StatusOK, info)
}

func (h *Handler) TTL(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	remaining, hasTTL, err := h.eng.TTL(key)
	if err != nil {
		code, msg := parseEngineError(err)
		sendError(w, code, msg)
		return
	}

	sendResponse(w, http.StatusOK, map[string]any{
		"key":       key,
		"has_ttl":   hasTTL,
		"ttl_ms":    remaining.Milliseconds(),
		"ttl_human": remaining.String(),
	})
}

func (h *Handler) Expire(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var body ExpireReq
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendError(w, http.StatusBadRequest, "invalid body")
		return
	}

	ttl, err := parseTTL(body.TTL)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateTTL(ttl); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.eng.Expire(key, ttl); err != nil {
		code, msg := parseEngineError(err)
		sendError(w, code, msg)
		return
	}

	sendResponse(w, http.StatusOK, nil)
}

func (h *Handler) Persist(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.eng.Persist(key); err != nil {
		code, msg := parseEngineError(err)
		sendError(w, code, msg)
		return
	}

	sendResponse(w, http.StatusOK, nil)
}
