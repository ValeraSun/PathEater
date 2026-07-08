package entity

import (
    "net"
    "sync"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID int64
    HP int
    X float32
    Y float32
}