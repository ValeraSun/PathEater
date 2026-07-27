package components

import "github.com/ValeraSun/PathEater/internal/core/types"

type RoomComponent struct {
	Players      map[types.Entity]bool
	ExternalWall types.Entity
	HasBreakdown bool
	Oxygen       float64
	Vacuum       bool
	MinX         float64
	MaxX         float64
	MinY         float64
	MaxY         float64
}

func (*RoomComponent) Type() string {
	return "room"
}

func NewRoomComponent(minX, maxX, minY, maxY float64, externalWall types.Entity) *RoomComponent {
	return &RoomComponent{
		Players:      make(map[types.Entity]bool),
		ExternalWall: externalWall,
		HasBreakdown: false,
		Oxygen:       100.0,
		Vacuum:       false,
		MinX:         minX,
		MaxX:         maxX,
		MinY:         minY,
		MaxY:         maxY,
	}
}
