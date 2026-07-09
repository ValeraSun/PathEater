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
    wg         sync.WaitGroup
    closed     bool
    closedMu   sync.RWMutex
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
    GameState map[string]interface{}
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
    h.wg.Add(1)
    defer h.wg.Done()
    log.Println("Hub запущен")

    for {
        select {
        case <-h.ctx.Done():
            log.Println("Hub остановлен")
            return

        case client, ok := <-h.register:
            if !ok {
                return
            }
            h.registerClient(client)

        case client, ok := <-h.unregister:
            if !ok {
                return
            }
            h.unregisterClient(client)

        case msg, ok := <-h.broadcast:
            if !ok {
                return
            }
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

    h.gamesMu.RLock()
    roomIDs := make([]string, 0, len(h.games))
    for id := range h.games {
        roomIDs = append(roomIDs, id)
    }
    h.gamesMu.RUnlock()

    for _, id := range roomIDs {
        h.gamesMu.RLock()
        room, exists := h.games[id]
        h.gamesMu.RUnlock()
        if !exists {
            continue
        }
        room.Mutex.Lock()
        delete(room.Clients, client.ID)
        if len(room.Clients) == 0 {
            room.Mutex.Unlock()
            h.DeleteGameRoom(id)
        } else {
            room.Mutex.Unlock()
        }
    }

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
    clients := make(map[string]*Client, len(h.clients))
    for id, client := range h.clients {
        clients[id] = client
    }
    h.clientsMu.RUnlock()

    for id, client := range clients {
        if excludeMap[id] {
            continue
        }
        if err := client.Send(data); err != nil {
            log.Printf("Ошибка отправки клиенту %s: %v", id, err)
            h.UnregisterClient(client)
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
    clients := make(map[string]*Client, len(room.Clients))
    for id, client := range room.Clients {
        clients[id] = client
    }
    room.Mutex.RUnlock()

    for id, client := range clients {
        if excludeMap[id] {
            continue
        }
        if err := client.Send(data); err != nil {
            log.Printf("Ошибка отправки клиенту %s в комнате %s: %v", id, roomID, err)
            h.UnregisterClient(client)
        }
    }
}

//Регистрирует клиента
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
    default:
        log.Printf("Канал регистрации переполнен для клиента %s", client.ID)
    }
}

//Удаляет клиента
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
    default:
        log.Printf("Канал удаления переполнен для клиента %s", client.ID)
    }
}

//Отправляет сообщение всем
func (h *Hub) Broadcast(msgType string, data interface{}, exclude ...string) {
    h.closedMu.RLock()
    if h.closed {
        h.closedMu.RUnlock()
        return
    }
    h.closedMu.RUnlock()

    select {
    case <-h.ctx.Done():
        return
    case h.broadcast <- BroadcastMessage{
        Type:    msgType,
        Data:    data,
        Exclude: exclude,
    }:
    default:
        log.Printf("Канал broadcast переполнен")
    }
}

//Отправляет сообщение в комнату
func (h *Hub) BroadcastToRoom(roomID string, msgType string, data interface{}, exclude ...string) {
    h.closedMu.RLock()
    if h.closed {
        h.closedMu.RUnlock()
        return
    }
    h.closedMu.RUnlock()

    select {
    case <-h.ctx.Done():
        return
    case h.broadcast <- BroadcastMessage{
        Type:    msgType,
        Data:    data,
        RoomID:  roomID,
        Exclude: exclude,
    }:
    default:
        log.Printf("Канал broadcast переполнен для комнаты %s", roomID)
    }
}

//Создает новую игровую комнату
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

//Возвращает игровую комнату
func (h *Hub) GetGameRoom(roomID string) (*GameRoom, bool) {
    h.gamesMu.RLock()
    defer h.gamesMu.RUnlock()
    room, exists := h.games[roomID]
    return room, exists
}

//Удаляет игровую комнату
func (h *Hub) DeleteGameRoom(roomID string) {
    h.gamesMu.Lock()
    room, exists := h.games[roomID]
    if exists {
        room.Mutex.Lock()
        for clientID := range room.Clients {
            delete(room.Clients, clientID)
        }
        room.Mutex.Unlock()
    }
    delete(h.games, roomID)
    h.gamesMu.Unlock()
    log.Printf("Удалена игровая комната %s", roomID)
}

//Возвращает количество клиентов в хабе
func (h *Hub) GetClientsCount() int {
    h.clientsMu.RLock()
    defer h.clientsMu.RUnlock()
    return len(h.clients)
}

//Возвращает количество комнат в хабе
func (h *Hub) GetGameRoomsCount() int {
    h.gamesMu.RLock()
    defer h.gamesMu.RUnlock()
    return len(h.games)
}

//Останавливает Hub
func (h *Hub) Shutdown() {
    h.closedMu.Lock()
    if h.closed {
        h.closedMu.Unlock()
        return
    }
    h.closed = true
    h.closedMu.Unlock()

    h.cancel()

    h.clientsMu.Lock()
    clients := make(map[string]*Client, len(h.clients))
    for id, client := range h.clients {
        clients[id] = client
    }
    h.clientsMu.Unlock()

    for _, client := range clients {
        client.Close()
    }

    close(h.register)
    close(h.unregister)
    close(h.broadcast)

    h.wg.Wait()

    log.Println("Hub полностью остановлен")
}

//Добавляет клиента в комнату
func (r *GameRoom) AddClient(client *Client) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()
    r.Clients[client.ID] = client
    log.Printf("Клиент %s добавлен в комнату %s", client.ID, r.ID)
}

//Удаляет клиента из комнаты
func (r *GameRoom) RemoveClient(clientID string) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()
    delete(r.Clients, clientID)
    log.Printf("Клиент %s удален из комнаты %s", clientID, r.ID)

    if len(r.Clients) == 0 {
        r.Hub.DeleteGameRoom(r.ID)
    }
}

//Отправляет сообщение в комнату
func (r *GameRoom) Broadcast(msgType string, data interface{}, exclude ...string) {
    r.Hub.BroadcastToRoom(r.ID, msgType, data, exclude...)
}

//Возвращает список клиентов в комнате
func (r *GameRoom) GetClients() []*Client {
    r.Mutex.RLock()
    defer r.Mutex.RUnlock()

    clients := make([]*Client, 0, len(r.Clients))
    for _, client := range r.Clients {
        clients = append(clients, client)
    }
    return clients
}

//Возвращает количество клиентов в комнате
func (r *GameRoom) GetClientsCount() int {
    r.Mutex.RLock()
    defer r.Mutex.RUnlock()
    return len(r.Clients)
}

//Получить состояние игры
func (r *GameRoom) GetGameState() map[string]interface{} {
    r.Mutex.RLock()
    defer r.Mutex.RUnlock()
    state := make(map[string]interface{})
    for k, v := range r.GameState {
        state[k] = v
    }
    return state
}

//Задать игре состояние
func (r *GameRoom) SetGameState(key string, value interface{}) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()
    r.GameState[key] = value
}

//Удалить состояние игры
func (r *GameRoom) DeleteGameState(key string) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()
    delete(r.GameState, key)
}

//Очистить состояние игры
func (r *GameRoom) ClearGameState() {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()
    r.GameState = make(map[string]interface{})
}