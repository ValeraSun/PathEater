package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // разрешает подключения с любого origin (только для разработки)
	},
}

type GameServer struct {
	clients   map[*websocket.Conn]*Player
	broadcast chan []byte
	mutex     sync.Mutex
}


type Player struct {
	ID      string
	X, Y, Z float64
}

// NewGameServer создаёт сервер
func NewGameServer() *GameServer {
	return &GameServer{
		clients:   make(map[*websocket.Conn]*Player),
		broadcast: make(chan []byte),
	}
}

func (server *GameServer) gameLoop(tickDuration time.Duration) {
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	for range ticker.C {
		server.Lock()
		//место для обработки запросов
		server.Unlock()
	}
}

// handleConnection - метод сервера, подключающий к нему нового игрока
func (server *GameServer) handleConnection(webSockConn *websocket.Conn) {
	defer func() {
		server.mutex.Lock()
		delete(server.clients, webSockConn)
		server.mutex.Unlock()
		webSockConn.Close()
	}()
	
	server.mutex.Lock()

	// Генерируем уникальный ID (для теста подойдёт простой счётчик или UUID)
	playerID := fmt.Sprintf("player_%d", time.Now().UnixNano())

	// Создаём игрока со стартовыми координатами
	player := NewPlayer(playerID)
	s.clients[webSockConn] = player

	// Сигнализируем, что игрок создался
	log.Printf("player %s created at (0,0,0)", player.ID)

	s.mutex.Unlock()

	for {
		_, message, err := webSockConn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		server.broadcast <- message
	}
}

func NewPlayer(playerID string) *Player{
	return $Player {
		ID: playerID,
		X:  0.0,
		Y:  0.0,
		Z:  0.0,
	}
}

//Место для основной обработки данных с клиентов
//Пока что пустышка, отправляющая в консоль код ошибки, при наличии таковой
func (server *GameServer) broadcaster() {
	for message := range server.broadcast {
		server.mutex.Lock()
		for client := range server.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("write error:", err)
				client.Close()
			}
		}
		server.mutex.Unlock()
	}
}

//handlerWebSocket преобразует HTTP-запрос в WebSocket-связь и, если всё удачно, добавляет на сервер нового игрока
func (server *GameServer) handlerWebSocket(respWriter http.ResponseWriter, request *http.Request) {
	webSockConn, err := upgrader.Upgrade(respWriter, request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	go server.handleConnection(webSockConn)
}

func main() {
	server := NewGameServer()

	//Место для основной обработки данных клиентов
	go server.broadcaster()

	http.HandleFunc("/ws", server.handlerWebSocket)

	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}