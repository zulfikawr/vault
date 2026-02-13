package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/realtime"
)

type RealtimeHandler struct {
	hub *realtime.Hub
}

func NewRealtimeHandler(hub *realtime.Hub) *RealtimeHandler {
	return &RealtimeHandler{hub: hub}
}

func (h *RealtimeHandler) Connect(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := make(realtime.Client, 10)
	h.hub.Register(client)

	defer h.hub.Unregister(client)

	// Keep-alive heartbeat
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			if _, err := fmt.Fprintf(w, ": ping\n\n"); err != nil {
				return // Client disconnected
			}
			flusher.Flush()
		case msg := <-client:
			if msg == nil {
				return
			}
			data, err := json.Marshal(msg)
			if err != nil {
				errors.Log(r.Context(), err, "marshal SSE message")
				continue
			}
			if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg.Action, string(data)); err != nil {
				return // Client disconnected
			}
			flusher.Flush()
		}
	}
}
