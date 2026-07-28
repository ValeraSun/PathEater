package events

type VacuumRecalculateEvent struct{}

func (*VacuumRecalculateEvent) Type() string { return "vacuumRecalculate" }

func NewVacuumRecalculateEvent() *VacuumRecalculateEvent {
	return &VacuumRecalculateEvent{}
}
