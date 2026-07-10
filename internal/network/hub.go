package network

import (
    "sync"
    "log"
    "context"
)

//структура хаба
type Hub struct {
    Clients    map[string]*Client
    clientsMu  sync.RWMutex

    Rooms      map[string]*GameRoom
    roomsMu    sync.RWMutex

    register   chan *Client
    unregister chan *Client

    wg         sync.WaitGroup

    ctx        context.Context
    cancel     context.CancelFunc

    closed     bool
    closedMu   sync.RWMutex
}

//создание хаба
var (
    hubInstance *Hub
    hubOnce     sync.Once
)
//Возвращает хаб
func GetHub() *Hub {
    hubOnce.Do(func() {
        hubInstance = newHub()
        go hubInstance.run()
    })
    return hubInstance
}

//Инициализирует хаб
func newHub() *Hub {
    ctx, cancel := context.WithCancel(context.Background())
    return &Hub{
    Clients:    make(map[string]*Client),
    Rooms:      make(map[string]*GameRoom),
    register:   make(chan *Client, 100),
    unregister: make(chan *Client, 100),
    ctx:        ctx,
    cancel:     cancel,
    closed:     false,
    }
}

func (h *Hub) run() {
    h.wg.Add(1)
    defer h.wg.Done()

    for {
        select {
        case <- h.ctx.Done():
            return
        case client, ok := <- h.register:
            if !ok {
                return
            }
            h.registerClient(client)
        
        case client, ok := <- h.unregister:
            if !ok {
                return
            }
            h.unregisterClient(client)
        }
    }
}

func (h *Hub) registerClient(client *Client) {
    h.clientsMu.Lock()
    h.Clients[client.ID] = client
    h.clientsMu.Unlock()
}

func (h *Hub) unregisterClient(client *Client) {
    h.clientsMu.Lock()
    delete(h.Clients, client.ID)
    h.clientsMu.Unlock()

    if client.Room != nil {
        client.Room.Mutex.Lock()
        delete(client.Room.Clients, client.ID)
        isEmpty := len(client.Room.Clients) == 0
        roomID := client.Room.ID
        client.Room.Mutex.Unlock()

        if isEmpty {
            h.Rooms[roomID].Close()
        }
    }
}

func (h *Hub) RegisterClient(client *Client) {
    h.closedMu.RLock()
    if h.closed {
        h.closedMu.RUnlock()
        return
    }
    h.closedMu.RUnlock()

    select {
    case <-h.ctx.Done():
        return
    case h.register <- client:
        return
    default:
        log.Printf("Канал регистрации переполнен для клиента %s", client.ID)
    }
}

func (h *Hub) UnregisterClient(client *Client) {
    h.closedMu.RLock()
    if h.closed {
        h.closedMu.RUnlock()
        return
    }
    h.closedMu.RUnlock()

    select {
    case <-h.ctx.Done():
        return
    case h.unregister <- client:
        return
    default:
        log.Printf("Канал удаления переполнен для клиента %s", client.ID)
    }
}

func (h *Hub) Close() {
    h.closedMu.Lock()
    if h.closed {
        h.closedMu.Unlock()
        return
    }
    h.closed = true
    h.closedMu.Unlock()

    h.roomsMu.RLock()
    rooms := make(map[string]*GameRoom, len(h.Rooms))
    for id, room := range h.Rooms {
        rooms[id] = room
    } 
    h.roomsMu.RUnlock()
    
    h.clientsMu.RLock()
    clients := make(map[string]*Client, len(h.Clients))
    for id, client := range h.Clients {
        clients[id] = client
    } 
    h.clientsMu.RUnlock()

    for _, client := range clients {
        client.Close()
    }

    for _, room := range rooms {
        room.Close()
    }

    close(h.register)
    close(h.unregister)

    h.cancel()

    h.wg.Wait()
}

//создаёт новую игровую комнату
func (h *Hub) NewGameRoom(roomID string) *GameRoom {
    h.roomsMu.Lock()
    defer h.roomsMu.Unlock()

    room := &GameRoom{
        ID:      roomID,
        Clients: make(map[string]*Client),
		World    NewWorld(),
    }

    h.Rooms[roomID] = room

    return room
}

//Возвращает игровую комнату
func (h *Hub) GetGameRoom(roomID string) (*GameRoom, bool) {
    h.roomsMu.RLock()
    defer h.roomsMu.RUnlock()
    room, exists := h.Rooms[roomID]
    return room, exists
}