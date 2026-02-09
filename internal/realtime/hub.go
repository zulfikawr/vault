package realtime

import (
	"context"
	"log/slog"
	"sync"
)

type Client chan *Message

type Hub struct {
	clients    map[Client]bool
	broadcast  chan *Message
	register   chan Client
	unregister chan Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan Client),
		unregister: make(chan Client),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			slog.Debug("Realtime client registered")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client)
			}
			h.mu.Unlock()
			slog.Debug("Realtime client unregistered")
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client <- message:
				default:
					// If client buffer is full, we might want to skip or disconnect
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(msg *Message) {
	h.broadcast <- msg
}

func (h *Hub) Register(c Client) {
	h.register <- c
}

func (h *Hub) Unregister(c Client) {
	h.unregister <- c
}
