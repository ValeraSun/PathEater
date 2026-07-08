package usecase

import (
	"log"
	"sync"
)

type Room struct {
	clients map[*Client]bool
	broadcast chan []byte
	closed bool
	mutex sync.RWMutex
}

//создание комнаты
func NewRoom() *Room {
	return &Room{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		closed: false,
	}
}

//закрытие комнаты
func (r *Room) Close() {
    r.mutex.Lock()
	defer r.mutex.Unlock()

    if r.closed {
        return
    }
    r.closed = true

    close(r.broadcast)

    for client := range r.clients {
        client.Close()
    }

    r.clients = nil
}

//рассылка сообщения клиентам комнаты
func (r *Room) Broadcaster() {
    for msg := range r.broadcast {
        r.mutex.RLock()
        for client := range r.clients {
            select {
            case client.send <- msg:
            default:
                go client.Close()
                delete(r.clients, client)
            }
        }
        r.mutex.RUnlock()
    }
}