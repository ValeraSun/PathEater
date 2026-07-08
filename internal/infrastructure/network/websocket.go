package network

import (
    "encoding/json"
	"log"
	"net"
	"net/http"
    "time"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true}, //ИЗМЕНИТЬ!!!
}

func RegisterHandlers() {
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		HandleConnection(writer, request)
	})
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка установления ws-связи:", err)
		return
	}
	defer wsConn.Close()
	

    udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8081") // Замените на ваш адрес
    if err != nil {
        log.Println("Ошибка резолва UDP адреса:", err)
        wsConn.Close()
        return
    }

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        log.Println("Ошибка подключения UDP:", err)
        return
    }
    defer udpConn.Close()


	client := NewClient(wsConn, udpConn, udpAddr)

	client.Start()
	
    go readMessages(wsConn, udpConn)

	err = wsConn.WriteMessage(websocket.TextMessage, []byte("Всё клёво!"))
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
		return
	}

	log.Println("Всё клёво!")

    <-client.done
    log.Println("Соединение закрыто клиентом")
}

func readMessages(wsConn *websocket.Conn, udpConn *net.UDPConn) {
    var messageCount int
    
    for {
        _, message, err := wsConn.ReadMessage()
        if err != nil {
            if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
                log.Println("Клиент закрыл соединение")
            } else {
                log.Println("Ошибка чтения сообщения:", err)
            }
            return
        } 
        
        messageCount++
        log.Printf("📨 [%d] Получено сообщение: %s", messageCount, string(message))
        
        // Отправляем в UDP
        _, err = udpConn.Write(message)
        if err != nil {
            log.Printf("Ошибка отправки в UDP: %v", err)
            continue
        }
        
        // Отправляем подтверждение клиенту
        ack := map[string]interface{}{
            "status": "ok",
            "received": messageCount,
            "timestamp": time.Now().Unix(),
        }
        ackData, _ := json.Marshal(ack)
        err = wsConn.WriteMessage(websocket.TextMessage, ackData)
        if err != nil {
            log.Printf("Ошибка отправки подтверждения: %v", err)
        }
    }
}

func handleMessage(message []byte, udpConn *net.UDPConn) {
    log.Printf("Получено сообщение: %s", string(message))
    
    // Здесь логика обработки сообщения
    // Например, отправка в UDP
    _, err := udpConn.Write(message)
    if err != nil {
        log.Println("Ошибка отправки в UDP:", err)
    }
}