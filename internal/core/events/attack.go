package events

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type Target int

const (
	Enemy Target = iota
	Friend
)

type AttackEvent struct {
	Damage   int
	Collider geometry.Collider
	Team     Target
}

func (*AttackEvent) Type() string { return "attack" }

func NewAttackEvent(damage int, collider geometry.Collider, target Target) *AttackEvent {
	return &AttackEvent{
		Damage:   damage,
		Collider: collider,
		Team:     target,
	}
}
