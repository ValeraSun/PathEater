package events

import "github.com/ValeraSun/PathEater/internal/core/types"

type BoardingEvent struct {
	ID types.Entity
}

func (*BoardingEvent) Type() string { return "boarding" }

func NewBoardingEvent(id types.Entity) *BoardingEvent {
	return &BoardingEvent{
		ID: id,
	}
}
