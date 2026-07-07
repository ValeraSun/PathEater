package room

import "github.com/ValeraSun/PathEater/internal/game/player"

type Room struct {
	ID       string
	Sessions map[string]*player.Session
	Started  bool
}

func (r *Room) AddSession(s *player.Session) {
	r.Sessions[s.Player.ID] = s
}

func (r *Room) BroadcastCommand(cmd *view.RenderCommand) {
	for _, s := range r.Sessions {
		_ = s.SendCommand(cmd)
	}
}