package components

import "github.com/ValeraSun/PathEater/internal/core/types"

type DoorComponent struct {
	RoomA  types.Entity
	RoomB  types.Entity
	IsOpen bool
}

func (*DoorComponent) Type() string {
	return "door"
}

func NewDoorComponent(roomA, roomB types.Entity) *DoorComponent {
	return &DoorComponent{
		RoomA:  roomA,
		RoomB:  roomB,
		IsOpen: false,
	}
}
