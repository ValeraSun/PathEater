package network

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/ValeraSun/PathEater/internal/core/ecs"
)

type Hub struct {
	register   chan *Client
	unregister chan *Client

	clients    map[string]*Client
	clientsMu  sync.RWMutex

	rooms      map[string]*GameRoom
	roomsMu    sync.RWMutex

	ctx        context.Context
	cancel     context.CancelFunc

	wg         sync.WaitGroup

	closed     bool
	closedMu   sync.RWMutex
}

type GameRoom struct {
	ID        string
	Clients   map[string]*Client
	Mutex     sync.RWMutex
	Hub       *Hub
	World     *ecs.World
}

var (
	hubInstance *Hub
	hubOnce     sync.Once
)

// Возвращает синглтон Hub
func GetHub() *Hub {
	hubOnce.Do(func() {
		hubInstance = NewHub()
		go hubInstance.Run()
	})
	return hubInstance
}

// Создает новый экземпляр Hub
func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		clients:    make(map[string]*Client),
		rooms:      make(map[string]*GameRoom),
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
		}
	}
}

// Регистрирует нового клиента в хабе
func (h *Hub) registerClient(client *Client) {
	h.clientsMu.Lock()
	h.clients[client.ID] = client
	h.clientsMu.Unlock()
}

// Удаляет клиента из хаба
func (h *Hub) unregisterClient(client *Client) {
	h.clientsMu.Lock()
	delete(h.clients, client.ID)
	h.clientsMu.Unlock()

	h.roomsMu.RLock()
	roomIDs := make([]string, 0, len(h.rooms))
	for id := range h.rooms {
		roomIDs = append(roomIDs, id)
	}
	h.roomsMu.RUnlock()

	for _, id := range roomIDs {
		h.roomsMu.RLock()
		room, exists := h.rooms[id]
		h.roomsMu.RUnlock()
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
}

// Регистрирует клиента
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

// Удаляет клиента
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

// Создает новую игровую комнату
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

// Возвращает игровую комнату
func (h *Hub) GetGameRoom(roomID string) (*GameRoom, bool) {
	h.gamesMu.RLock()
	defer h.gamesMu.RUnlock()
	room, exists := h.games[roomID]
	return room, exists
}

// Возвращает количество клиентов в хабе
func (h *Hub) GetClientsCount() int {
	h.clientsMu.RLock()
	defer h.clientsMu.RUnlock()
	return len(h.clients)
}

// Возвращает количество комнат в хабе
func (h *Hub) GetGameRoomsCount() int {
	h.gamesMu.RLock()
	defer h.gamesMu.RUnlock()
	return len(h.games)
}

// Останавливает Hub
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

	h.wg.Wait()
}