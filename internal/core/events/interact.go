package events

import "github.com/ValeraSun/PathEater/internal/core/types"

type InteractTerminalEvent struct {
	ID types.Entity
}

func (*InteractTerminalEvent) Type() string { return "interactTerminal" }

func NewInteractTerminalEvent(id types.Entity) *InteractTerminalEvent {
	return &InteractTerminalEvent{
		ID: id,
	}
}
