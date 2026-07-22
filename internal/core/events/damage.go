package events

import "github.com/ValeraSun/PathEater/internal/core/types"

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
	ID     string
	Damage int
}

func (*DamageCosmoAlienEvent) Type() string { return "damageCosmoAlien" }

func NewDamageCosmoAlienEvent(id string, damage int) *DamageCosmoAlienEvent {
	return &DamageCosmoAlienEvent{
		ID:     id,
		Damage: damage,
	}
}
