package network

import (
	"log"
	//"time"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // разрешает подключения с любого origin (только для разработки)
	},
}

// HandleClient - метод сервера, подключающий к нему нового игрока
func HandleClient(webSock *websocket.Conn, hub *Hub) {
	defer webSock.Close()

	hub.AddClient(webSock)
	defer hub.RemoveClient(webSock)

	/*webSock.SetPongHandler(func(message string) error {
		// опционально: можно логировать или обновлять lastPongTime у игрока
		return nil
	})
	pingTicker := time.NewTicker(15 * time.Second)
	defer pingTicker.Stop()*/ //блок для проверки, жив ли клиент

	for {
		mt, message, err := webSock.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			break
		}

		if mt != websocket.TextMessage {
			continue
		}

		hub.SendToAll(message) // или можно слать в отдельный канал команд для game.go
	}
}