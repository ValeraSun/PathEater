package network

import (
	"sync"
)

type GameRoom struct {
    ID        string
    Clients   map[string]*Client
    Mutex     sync.RWMutex
    World     World
}

//добавляет клиента в комнату
func (r *GameRoom) AddClient(client *Client) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Clients[client.ID] = client
	client.Room = r
}

//удаляет клиента из комнаты
func (r *GameRoom) RemoveClient(client *Client) {
	r.Mutex.Lock()
	delete(r.Clients, client.ID)
	client.Room = nil
    isEmpty := len(r.Clients) == 0
    r.Mutex.Unlock()

    if isEmpty {
        r.Close()
    }
}


//закрывает комнату
func (r *GameRoom) Close() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

    for _, client := range r.Clients {
	    client.mu.Lock()
        client.Room = nil
    }
    r.Clients = make(map[string]*Client)
	
    if r.World != nil {
        r.World.Close()
    }

	r.ID = ""
}