package events

import "github.com/ValeraSun/PathEater/internal/core/types"

type AnimationEvent struct {
	State    string
	ID       types.Entity
	Cooldown float32
}

func (*AnimationEvent) Type() string { return "animation" }

func NewAnimationEvent(id types.Entity, state string, cooldwon float32) *AnimationEvent {
	return &AnimationEvent{
		state,
		id,
		cooldwon,
	}
}
