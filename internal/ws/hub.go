package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu    sync.Mutex
	conns map[string]map[*websocket.Conn]struct{}
	upg   websocket.Upgrader
}

func NewHub() *Hub {
	return &Hub{
		conns: make(map[string]map[*websocket.Conn]struct{}),
		upg:   websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		topic = "default"
	}
	c, err := h.upg.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.mu.Lock()
	if h.conns[topic] == nil {
		h.conns[topic] = make(map[*websocket.Conn]struct{})
	}
	h.conns[topic][c] = struct{}{}
	h.mu.Unlock()

	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.conns[topic], c)
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
func (h *Hub) Broadcast(topic string, v interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if topic == "" {
		for _, conns := range h.conns {
			for c := range conns {
				_ = c.WriteJSON(v)
			}
		}
		return
	}
	for c := range h.conns[topic] {
		_ = c.WriteJSON(v)
	}
}
