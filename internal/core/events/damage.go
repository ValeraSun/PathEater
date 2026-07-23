package events

import "github.com/ValeraSun/PathEater/internal/core/types"

type DamageDealEvent struct {
	ID     types.Entity
	Damage int
}

func (*DamageDealEvent) Type() string { return "damageDeal" }

func NewDamageDealEvent(id types.Entity, damage int) *DamageDealEvent {
	return &DamageDealEvent{
		ID:     id,
		Damage: damage,
	}

}
