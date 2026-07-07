package network

import (
	"log"
	"net/http"
)

func RegisterHandlers(hub *Hub) {
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		handleWebSocketUpgrade(writer, request, hub)
	})
}

func handleWebSocketUpgrade(writer http.ResponseWriter, request *http.Request, hub *Hub) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ws, err := Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	go HandleClient(ws, hub)
}