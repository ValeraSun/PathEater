package usecase

import (
	"fmt"
)

//добавляет клиента в комнату
func (r *Room) AddClient(client *Client) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if r.closed {
        return fmt.Errorf("room is closed")
    }

    r.clients[client] = true
    return nil
}

//удаляет клиента из комнаты
func (r *Room) RemoveClient(client *Client) {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if r.closed {
        return
    }

    if _, exists := r.clients[client]; !exists {
        return // клиент уже не в комнате
    }

    delete(r.clients, client)
    client.Close() // закрываем соединение клиента
    
    log.Printf("Client removed from room. Remaining: %d", len(r.clients))
}