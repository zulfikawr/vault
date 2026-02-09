package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
			fmt.Fprintf(w, ": ping\n\n")
			flusher.Flush()
		case msg := <-client:
			if msg == nil {
				return
			}
			data, _ := json.Marshal(msg)
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg.Action, string(data))
			flusher.Flush()
		}
	}
}