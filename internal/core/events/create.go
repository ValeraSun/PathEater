package events

import "github.com/ValeraSun/PathEater/internal/core/geometry"

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

type CreateAlienEvent struct {
	Position geometry.Vec3
}

func (*CreateAlienEvent) Type() string { return "createAlien" }

func NewCreateAlienEvent(position geometry.Vec3) *CreateAlienEvent {
	return &CreateAlienEvent{
		Position: position,
	}
}
