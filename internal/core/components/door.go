package components

import "github.com/ValeraSun/PathEater/internal/core/types"

type DoorComponent struct {
	RoomA       types.Entity
	RoomB       types.Entity
	IsOpen      bool
	TriggerMinX float64
	TriggerMaxX float64
	TriggerMinY float64
	TriggerMaxY float64
}

func (*DoorComponent) Type() string {
	return "door"
}

func NewDoorComponent(roomA, roomB types.Entity, tminx, tmaxx, tminy, tmaxy float64) *DoorComponent {
	return &DoorComponent{
		RoomA:       roomA,
		RoomB:       roomB,
		IsOpen:      false,
		TriggerMinX: tminx,
		TriggerMaxX: tmaxx,
		TriggerMinY: tminy,
		TriggerMaxY: tmaxy,
	}
}
