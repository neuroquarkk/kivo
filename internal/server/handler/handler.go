package handler

import (
	"encoding/json"
	"net/http"
	"time"

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

	if err := validateTTL(time.Duration(body.TTL)); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.eng.Set(
		key, []byte(body.Value),
		time.Duration(body.TTL),
	); err != nil {
		code, msg := parseEngineError(err)
		sendError(w, code, msg)
		return
	}

	sendResponse(w, http.StatusCreated, nil)
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
