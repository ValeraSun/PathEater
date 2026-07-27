package events

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type DeleteEvent struct {
	ID types.Entity
}

func (*DeleteEvent) Type() string { return "delete" }

func NewDeleteEvent(id types.Entity) *DeleteEvent {
	return &DeleteEvent{
		ID: id,
	}
}
