package events

// локальные event через маленькую
// глобальные event через большую

type MeteoriteZoneEvent struct {
	Active bool
}

func (*MeteoriteZoneEvent) Type() string { return "meteoriteZone" }

func NewMeteoriteZoneEvent() *MeteoriteZoneEvent {
	return &MeteoriteZoneEvent{
		Active: false,
	}
}
