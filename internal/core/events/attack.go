package events

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type Target int

const (
	Enemy Target = iota
	Friend
)

type AttackEvent struct {
	Damage   int
	Collider geometry.Collider
	Team     Target
	Attacker types.Entity
}

func (*AttackEvent) Type() string { return "attack" }

func NewAttackEvent(damage int, collider geometry.Collider, attacker types.Entity) *AttackEvent {
	return &AttackEvent{
		Damage:   damage,
		Collider: collider,
		Attacker: attacker,
	}
}
