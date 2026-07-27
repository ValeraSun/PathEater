package events

type BreakdownEvent struct {
	SpawnAlien bool
}

func (*BreakdownEvent) Type() string { return "breakdown" }

func NewBreakdownEvent(spawnAlien bool) *BreakdownEvent {
	return &BreakdownEvent{
		SpawnAlien: spawnAlien,
	}
}
