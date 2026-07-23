package events

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type PlayerDeathEvent struct {
	ID types.Entity
}

func (*PlayerDeathEvent) Type() string { return "playerDeath" }

func NewPlayerDeathEvent(id types.Entity) *PlayerDeathEvent {
	return &PlayerDeathEvent{
		ID: id,
	}
}
