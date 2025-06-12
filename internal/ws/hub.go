package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]struct{}
	upg   websocket.Upgrader
}

func NewHub() *Hub {
	return &Hub{
		conns: make(map[*websocket.Conn]struct{}),
		upg:   websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	c, err := h.upg.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.mu.Lock()
	h.conns[c] = struct{}{}
	h.mu.Unlock()

	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.conns, c)
			h.mu.Unlock()
			c.Close()
		}()
		for {
			if _, _, err := c.NextReader(); err != nil {
				return
			}
		}
	}()
}

// Broadcast sends v to all connected clients.
func (h *Hub) Broadcast(v interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.conns {
		_ = c.WriteJSON(v)
	}
}
