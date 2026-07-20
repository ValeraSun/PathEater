package events

type DeathCosmoAlienEvent struct {
	ID     string
}

func (*DeathCosmoAlienEvent) Type() string { return "deathCosmoAlien" }

func NewDeathCosmoAlienEvent(id string) *DeathCosmoAlienEvent {
	return &DamageShipEvent{
		ID:     id,
	}
}