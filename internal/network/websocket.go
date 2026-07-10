package network

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true}, //ИЗМЕНИТЬ!!!
}

func RegisterHandlers() {
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		handleConnection(writer, request)
	})
	/*возможно, стоит добавить для реализации меню 
		http.HandleFunc("/http", func(writer http.ResponseWriter, request *http.Request) {
		HandleConnectionHttp(writer, request)
	})
	*/
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
    //Установка websocket-связи
	wsConn, err := upgrade.Upgrade(w, r, nil)
    //обработка ошибок
	if err != nil {
		log.Println("Ошибка установления ws-связи: ", err)
		return
	}
    defer wsConn.Close()
	
	//Инициализация/определение хаба
	hub := GetHub()

    //Создание клиента
	client := NewClient(wsConn)
    hub.RegisterClient(client)
    defer hub.UnregisterClient(client)
	
    //Запуск чтения сообщений от клиента и отправки сообщений от сервера
    go client.ReadMessages()
    go client.WriteMessages()

	//ожидание конца работы клиента
    <-client.done
}