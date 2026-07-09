package events

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

// локальные event через маленькую
// глобальные event через большую

type MoveEvent struct {
	Id        types.Entity
	Position  geometry.Vector3
	Direction geometry.Vector3
}

func (e *MoveEvent) Type() string { return "Move" }

func CreateEventMove(position geometry.Vector3, direction geometry.Vector3, id types.Entity) *MoveEvent {
	return &MoveEvent{
		Id:        id,
		Position:  position,
		Direction: direction,
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
