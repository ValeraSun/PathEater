package components

import "github.com/ValeraSun/PathEater/internal/core/types"

type PlayerComponent struct {
	RoomID types.Entity
	Dead   bool
}

func (*PlayerComponent) Type() string {
	return "player"
}

func NewPlayerComponent() *PlayerComponent {
	return &PlayerComponent{
		RoomID: "louge",
		Dead:   false,
	}
}
