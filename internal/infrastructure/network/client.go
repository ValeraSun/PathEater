package network

import (
    "log"
    "net"
    "sync"
    
    "github.com/gorilla/websocket"
)

type Client struct {
    wsConn   *websocket.Conn
    udpConn  *net.UDPConn
    udpAddr  *net.UDPAddr
	done     chan []byte
    mu       sync.Mutex
    isActive bool
}

func NewClient(wsConn *websocket.Conn, udpConn *net.UDPConn, udpAddr *net.UDPAddr) *Client {
    return &Client{
        wsConn:   wsConn,
        udpConn:  udpConn,
        udpAddr:  udpAddr,
		done:     make(chan []byte),
        isActive: true,
    }
}

func (c *Client) Start() {
    log.Println("Клиент запущен")
    // Здесь ваша логика старта клиента
}

func (c *Client) Close() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.isActive = false
    if c.wsConn != nil {
        c.wsConn.Close()
    }
    if c.udpConn != nil {
        c.udpConn.Close()
    }
}