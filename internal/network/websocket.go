package network

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // TODO: перед выкладкой ограничить своим доменом
}

// обрабатывает /ws запросы
func RegisterHandlers() {
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		HandleConnection(writer, request)
	})
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	hub := GetHub()

	wsConn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка установления ws-связи:", err)
		return
	}
	defer wsConn.Close()

	client := NewClient(wsConn)
	done := client.done

	hub.RegisterClient(client)
	defer hub.UnregisterClient(client)

	go client.ReadMessages()
	go client.WriteMessages()

	<-done
}