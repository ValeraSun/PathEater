package events

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type DoorEvent struct {
	ID types.Entity
}

func (*DoorEvent) Type() string { return "door" }

func NewDoorEvent(id types.Entity) *DoorEvent {
	return &DoorEvent{
		ID: id,
	}
}
