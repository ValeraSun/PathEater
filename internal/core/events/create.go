package events

type CreatePlayerEvent struct {
	ID string
}

func (*CreatePlayerEvent) Type() string { return "CreatePlayer" }

func NewCreatePlayerEvent(id string) *CreatePlayerEvent {
	return &CreatePlayerEvent{
		ID: id,
	}
}
