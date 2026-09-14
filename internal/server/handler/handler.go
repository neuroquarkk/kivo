package handler

import (
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
