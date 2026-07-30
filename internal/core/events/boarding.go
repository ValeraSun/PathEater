package events

type BoardingEvent struct{}

func (*BoardingEvent) Type() string { return "boarding" }

func NewBoardingEvent() *BoardingEvent {
	return &BoardingEvent{}
}
