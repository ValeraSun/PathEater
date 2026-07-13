package events

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

// локальные event через маленькую
// глобальные event через большую

type PlayerState struct {
	MoveFront bool             `json:"move_front"`
	MoveLeft  bool             `json:"move_left"`
	MoveRight bool             `json:"move_right"`
	MoveBack  bool             `json:"move_back"`
	Interact  bool             `json:"interact"`
	Attack    bool             `json:"attack"`
	Direction geometry.Vector3 `json:"direction"`
}

type SetPlayerStateEvent struct {
	PlayerState PlayerState
	ID          types.Entity
}

func (*SetPlayerStateEvent) Type() string { return "SetPlayerState" }
func NewSetPlayerStateEvent(ps PlayerState, id types.Entity) *SetPlayerStateEvent {
	return &SetPlayerStateEvent{
		PlayerState: ps,
		ID:          id,
	}
}

type EventUseItem struct {
	ItemID string
}

func (e *EventUseItem) Type() string { return "UseItem" }

func CreateEventUseItem(itemID string) EventUseItem {
	return EventUseItem{
		ItemID: itemID,
	}
}

type EventExit struct{}

func (e *EventExit) Type() string { return "Exit" }

func CreateEventExit() EventExit {
	return EventExit{}
}
