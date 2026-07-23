package events

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type DamagePlayerEvent struct {
	ID     types.Entity
	Damage int
}

func (*DamagePlayerEvent) Type() string { return "damagePlayer" }

func NewDamagePlayerEvent(id types.Entity, damage int) *DamagePlayerEvent {
	return &DamagePlayerEvent{
		ID:     id,
		Damage: damage,
	}
}

type DamageShipEvent struct {
	ID     types.Entity
	Damage int
}

func (*DamageShipEvent) Type() string { return "damageShip" }

func NewDamageShipEvent(id types.Entity, damage int) *DamageShipEvent {
	return &DamageShipEvent{
		ID:     id,
		Damage: damage,
	}
}

type DamageCosmoAlienEvent struct {
	ID     types.Entity
	Damage int
}

func (*DamageCosmoAlienEvent) Type() string { return "damageCosmoAlien" }

func NewDamageCosmoAlienEvent(id types.Entity, damage int) *DamageCosmoAlienEvent {
	return &DamageCosmoAlienEvent{
		ID:     id,
		Damage: damage,
	}
}
