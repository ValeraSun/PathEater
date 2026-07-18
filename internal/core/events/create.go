package events

type CreatePlayerEvent struct {
	ID string
}

func (*CreatePlayerEvent) Type() string { return "createPlayer" }

func NewCreatePlayerEvent(id string) *CreatePlayerEvent {
	return &CreatePlayerEvent{
		ID: id,
	}
}

type CreateShipEvent struct{}

func (*CreateShipEvent) Type() string { return "createShip" }
func NewCreateShipEvent() *CreateShipEvent {
	return &CreateShipEvent{}
}
