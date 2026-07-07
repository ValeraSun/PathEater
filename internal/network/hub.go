package network

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	register  chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (hub *Hub) AddClient(client *websocket.Conn) {
	hub.register <- client
}

func (hub *Hub) RemoveClient(client *websocket.Conn) {
	hub.unregister <- client
}

func (hub *Hub) SendToAll(data []byte) {
	hub.broadcast <- data
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.mutex.Lock()
			hub.clients[client] = true
			hub.mutex.Unlock()

		case client := <-hub.unregister:
			hub.mutex.Lock()
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				client.Close()
			}
			hub.mutex.Unlock()

		case message := <-hub.broadcast:
			hub.mutex.RLock()
			for client := range hub.clients {
				// можно добавить проверку на backpressure (если send-канал переполнен)
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					hub.unregister <- client
				}
			}
			hub.mutex.RUnlock()
		}
	}
}