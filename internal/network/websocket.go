package network

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait    = 60 * time.Second
	pingInterval = (pongWait * 9) / 10
)

func SetupPingPong(c *websocket.Conn) {
	c.SetPongHandler(func(string) error {
		// опционально: обновить lastPongTime у клиента
		return nil
	})
	// можно также запустить тикер для отправки ping из отдельной горутины
}