package network

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Hub  *Hub
}

func HandleClient(ws *websocket.Conn, h *Hub) {
	c := &Client{Conn: ws, Hub: h}

	h.register <- c
	defer func() {
		h.unregister <- c
	}()

	// опционально: настроить пинг/понг
	// websocketutil.SetupPingPong(ws)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		// здесь можно валидировать сообщение, парсить JSON и т.п.
		h.broadcast <- message
	}
}