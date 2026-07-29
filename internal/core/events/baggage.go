package events

type DamageBaggageEvent struct{}

func (*DamageBaggageEvent) Type() string { return "damageBaggage" }

func NewDamageBaggageEvent() *DamageBaggageEvent {
	return &DamageBaggageEvent{}
}
