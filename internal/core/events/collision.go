package events

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type CollisionEvent struct {
	ID1 types.Entity
	ID2 types.Entity
}

func (*CollisionEvent) Type() string { return "collision" }

func NewCollisionEvent(id1, id2 types.Entity) *CollisionEvent {
	return &CollisionEvent{
		ID1: id1,
		ID2: id2,
	}
}
