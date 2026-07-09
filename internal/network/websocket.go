package network

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //ИЗМЕНИТЬ!!!
}

func RegisterHandlers() {
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		HandleConnection(writer, request)
	})
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	hub := GetHub()
	//Установка websocket-связи
	wsConn, err := upgrade.Upgrade(w, r, nil)
	//обработка ошибок
	if err != nil {
		log.Println("Ошибка установления ws-связи:", err)
		return
	}
	defer wsConn.Close()

	//Создание клиента
	client := NewClient(wsConn)
	hub.RegisterClient(client)
	defer hub.UnregisterClient(client)

	//Запуск чтения сообщений от клиента и отправки сообщений от сервера
	go client.ReadMessages()
	go client.WriteMessages()

	<-client.done
}
