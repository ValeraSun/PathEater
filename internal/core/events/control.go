package events

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

// локальные event через маленькую
// глобальные event через большую

type PlayerState struct {
	MoveFront bool          `json:"move_front"`
	MoveLeft  bool          `json:"move_left"`
	MoveRight bool          `json:"move_right"`
	MoveBack  bool          `json:"move_back"`
	Interact  bool          `json:"interact"`
	Attack    bool          `json:"attack"`
	Direction geometry.Vec3 `json:"direction"`
}

type SetPlayerStateEvent struct {
	PlayerState PlayerState
	ID          string
}

func (*SetPlayerStateEvent) Type() string { return "setPlayerState" }

func NewSetPlayerStateEvent(ps PlayerState, id string) *SetPlayerStateEvent {
	return &SetPlayerStateEvent{
		PlayerState: ps,
		ID:          id,
	}
}
