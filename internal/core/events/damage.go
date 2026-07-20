package events

type DamageShipEvent struct {
	ID     string
	Damage int
}

func (*DamageShipEvent) Type() string { return "damageShip" }

func NewDamageShipEvent(id string, damage int) *DamageShipEvent {
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