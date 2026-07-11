package events

import (
	"log"

	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type MoveEvent struct {
	Id        types.Entity
	Position  geometry.Vector3
	Direction geometry.Vector3
}

func (e *MoveEvent) Type() string { return "Move" }

func CreateEventMove(position geometry.Vector3, direction geometry.Vector3, id types.Entity) *MoveEvent {
	log.Println("создан EventMove", id)
	return &MoveEvent{
		Id:        id,
		Position:  position,
		Direction: direction,
	}
}
