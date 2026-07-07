package network

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients       map[*websocket.Conn]bool
	broadcast     chan []byte
	register      chan *Client
	unregister    chan *Client
	mutex         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mutex.Lock()
			h.clients[c.Conn] = true
			h.mutex.Unlock()
		case c := <-h.unregister:
			h.mutex.Lock()
			delete(h.clients, c.Conn)
			h.mutex.Unlock()
			c.Conn.Close()
		case message := <-h.broadcast:
			h.mutex.RLock()
			for conn := range h.clients {
				conn.WriteMessage(websocket.TextMessage, message)
			}
			h.mutex.RUnlock()
		}
	}
}