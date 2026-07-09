package network

import (
    "context"
    "encoding/json"
    "log"
    "sync"
)

type Hub struct {
    register   chan *Client
    unregister chan *Client
    clients    map[string]*Client
    clientsMu  sync.RWMutex
    games      map[string]*GameRoom
    gamesMu    sync.RWMutex
    broadcast  chan BroadcastMessage
    ctx        context.Context
    cancel     context.CancelFunc
}

type BroadcastMessage struct {
    Type    string      `json:"type"`
    Data    interface{} `json:"data"`
    RoomID  string      `json:"roomId,omitempty"`
    Exclude []string    `json:"exclude,omitempty"`
}

type GameRoom struct {
    ID        string
    Clients   map[string]*Client
    Mutex     sync.RWMutex
    GameState interface{}
    Hub       *Hub
}

var (
    hubInstance *Hub
    hubOnce     sync.Once
)

//Возвращает синглтон Hub
func GetHub() *Hub {
    hubOnce.Do(func() {
        hubInstance = NewHub()
        go hubInstance.Run()
    })
    return hubInstance
}

//Создает новый экземпляр Hub
func NewHub() *Hub {
    ctx, cancel := context.WithCancel(context.Background())
    return &Hub{
        register:   make(chan *Client, 100),
        unregister: make(chan *Client, 100),
        clients:    make(map[string]*Client),
        games:      make(map[string]*GameRoom),
        broadcast:  make(chan BroadcastMessage, 1000),
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (h *Hub) Run() {
    log.Println("Hub запущен")
    
    for {
        select {
        case <-h.ctx.Done():
            log.Println("Hub остановлен")
            return
            
        case client := <-h.register:
            h.registerClient(client)
            
        case client := <-h.unregister:
            h.unregisterClient(client)
            
        case msg := <-h.broadcast:
            h.handleBroadcast(msg)
        }
    }
}

//Регистрирует нового клиента в хабе
func (h *Hub) registerClient(client *Client) {
    h.clientsMu.Lock()
    h.clients[client.ID] = client
    h.clientsMu.Unlock()
    
    log.Printf("Клиент %s подключен. Всего клиентов: %d", client.ID, len(h.clients))
}

//Удаляет клиента из хаба
func (h *Hub) unregisterClient(client *Client) {
    h.clientsMu.Lock()
    delete(h.clients, client.ID)
    h.clientsMu.Unlock()
    
    // Удаляем клиента из всех игр
    h.gamesMu.RLock()
    for _, game := range h.games {
        game.Mutex.Lock()
        delete(game.Clients, client.ID)
        game.Mutex.Unlock()
    }
    h.gamesMu.RUnlock()
    
    log.Printf("Клиент %s отключен. Всего клиентов: %d", client.ID, len(h.clients))
}

//Обрабатывает сообщения сервера
func (h *Hub) handleBroadcast(msg BroadcastMessage) {
    jsonData, err := json.Marshal(msg)
    if err != nil {
        log.Printf("Ошибка маршалинга broadcast: %v", err)
        return
    }
    
    if msg.RoomID != "" {
        h.sendToRoom(msg.RoomID, jsonData, msg.Exclude)
        return
    }
    
    h.sendToAll(jsonData, msg.Exclude)
}

//Отправляет сообщение всем клиентам
func (h *Hub) sendToAll(data []byte, exclude []string) {
    excludeMap := make(map[string]bool)
    for _, id := range exclude {
        excludeMap[id] = true
    }
    
    h.clientsMu.RLock()
    defer h.clientsMu.RUnlock()
    
    for id, client := range h.clients {
        if excludeMap[id] {
            continue
        }
        if err := client.Send(data); err != nil {
            log.Printf("Ошибка отправки клиенту %s: %v", id, err)
        }
    }
}

//Отправляет сообщение в комнату
func (h *Hub) sendToRoom(roomID string, data []byte, exclude []string) {
    h.gamesMu.RLock()
    room, exists := h.games[roomID]
    h.gamesMu.RUnlock()
    
    if !exists {
        log.Printf("Комната %s не найдена", roomID)
        return
    }
    
    excludeMap := make(map[string]bool)
    for _, id := range exclude {
        excludeMap[id] = true
    }
    
    room.Mutex.RLock()
    defer room.Mutex.RUnlock()
    
    for id, client := range room.Clients {
        if excludeMap[id] {
            continue
        }
        if err := client.Send(data); err != nil {
            log.Printf("Ошибка отправки клиенту %s в комнате %s: %v", id, roomID, err)
        }
    }
}

//Регистрирует клиента
func (h *Hub) RegisterClient(client *Client) {
    select {
    case h.register <- client:
    default:
        log.Printf("Канал регистрации переполнен для клиента %s", client.ID)
    }
}

//Удаляет клиента
func (h *Hub) UnregisterClient(client *Client) {
    select {
    case h.unregister <- client:
    default:
        log.Printf("Канал удаления переполнен для клиента %s", client.ID)
    }
}

// Broadcast отправляет сообщение всем
func (h *Hub) Broadcast(msgType string, data interface{}, exclude ...string) {
    h.broadcast <- BroadcastMessage{
        Type:    msgType,
        Data:    data,
        Exclude: exclude,
    }
}

// BroadcastToRoom отправляет сообщение в комнату
func (h *Hub) BroadcastToRoom(roomID string, msgType string, data interface{}, exclude ...string) {
    h.broadcast <- BroadcastMessage{
        Type:    msgType,
        Data:    data,
        RoomID:  roomID,
        Exclude: exclude,
    }
}

// CreateGameRoom создает новую игровую комнату
func (h *Hub) CreateGameRoom(roomID string) *GameRoom {
    h.gamesMu.Lock()
    defer h.gamesMu.Unlock()
    
    room := &GameRoom{
        ID:        roomID,
        Clients:   make(map[string]*Client),
        Hub:       h,
        GameState: make(map[string]interface{}),
    }
    
    h.games[roomID] = room
    log.Printf("Создана игровая комната %s", roomID)
    return room
}

// GetGameRoom возвращает игровую комнату
func (h *Hub) GetGameRoom(roomID string) (*GameRoom, bool) {
    h.gamesMu.RLock()
    defer h.gamesMu.RUnlock()
    room, exists := h.games[roomID]
    return room, exists
}

// DeleteGameRoom удаляет игровую комнату
func (h *Hub) DeleteGameRoom(roomID string) {
    h.gamesMu.Lock()
    delete(h.games, roomID)
    h.gamesMu.Unlock()
    log.Printf("Удалена игровая комната %s", roomID)
}

// GetClientsCount возвращает количество клиентов
func (h *Hub) GetClientsCount() int {
    h.clientsMu.RLock()
    defer h.clientsMu.RUnlock()
    return len(h.clients)
}

// GetGameRoomsCount возвращает количество комнат
func (h *Hub) GetGameRoomsCount() int {
    h.gamesMu.RLock()
    defer h.gamesMu.RUnlock()
    return len(h.games)
}

// Shutdown останавливает Hub
func (h *Hub) Shutdown() {
    h.cancel()
}

// Методы для GameRoom

// AddClient добавляет клиента в комнату
func (r *GameRoom) AddClient(client *Client) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()
    r.Clients[client.ID] = client
    log.Printf("Клиент %s добавлен в комнату %s", client.ID, r.ID)
}

// RemoveClient удаляет клиента из комнаты
func (r *GameRoom) RemoveClient(clientID string) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()
    delete(r.Clients, clientID)
    log.Printf("Клиент %s удален из комнаты %s", clientID, r.ID)
    
    // Если комната пуста - удаляем ее
    if len(r.Clients) == 0 {
        r.Hub.DeleteGameRoom(r.ID)
    }
}

// Broadcast отправляет сообщение в комнату
func (r *GameRoom) Broadcast(msgType string, data interface{}, exclude ...string) {
    r.Hub.BroadcastToRoom(r.ID, msgType, data, exclude...)
}

// GetClients возвращает список клиентов в комнате
func (r *GameRoom) GetClients() []*Client {
    r.Mutex.RLock()
    defer r.Mutex.RUnlock()
    
    clients := make([]*Client, 0, len(r.Clients))
    for _, client := range r.Clients {
        clients = append(clients, client)
    }
    return clients
}

// GetClientsCount возвращает количество клиентов в комнате
func (r *GameRoom) GetClientsCount() int {
    r.Mutex.RLock()
    defer r.Mutex.RUnlock()
    return len(r.Clients)
}