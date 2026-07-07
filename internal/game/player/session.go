package player

import (
	"github.com/gorilla/websocket"
	"github.com/ValeraSun/PathEater/internal/game/view"
)

type Session struct {
	Conn   *websocket.Conn
	Player *Player
	Ready  bool
}

func (s *Session) SendCommand(cmd *view.RenderCommand) error {
	// сериализация и отправка через ws
	return s.Conn.WriteJSON(cmd)
}