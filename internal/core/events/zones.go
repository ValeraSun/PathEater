package events

type MeteoriteZoneEvent struct{}

func (*MeteoriteZoneEvent) Type() string { return "meteoriteZone" }

func NewMeteoriteZoneEvent() *MeteoriteZoneEvent {
	return &MeteoriteZoneEvent{}
}
