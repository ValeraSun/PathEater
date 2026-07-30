package events

import "github.com/ValeraSun/PathEater/internal/core/types"

type DamageDealEvent struct {
	ID     types.Entity
	Damage float32
}

func (*DamageDealEvent) Type() string { return "damageDeal" }

func NewDamageDealEvent(id types.Entity, damage float32) *DamageDealEvent {
	return &DamageDealEvent{
		ID:     id,
		Damage: damage,
	}

}
