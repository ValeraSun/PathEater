package events

import "github.com/ValeraSun/PathEater/internal/core/types"

type DeadEvent struct {
	ID types.Entity
}

func (*DeadEvent) Type() string { return "dead" }

func NewDeadEvent(id types.Entity) *DeadEvent {

	return &DeadEvent{
		ID: id,
	}
}
