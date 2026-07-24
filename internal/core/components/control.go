package components

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type ControlComponent struct {
	MoveFront bool
	MoveLeft  bool
	MoveRight bool
	MoveBack  bool
	Interact  bool
	Attack    bool
	Direction geometry.Vec3
	ClientID  string
}

func (*ControlComponent) Type() string {
	return "control"
}

func NewControlComponent(clientID string) *ControlComponent {
	return &ControlComponent{
		ClientID: clientID,
	}
}

func (c *ControlComponent) Reset() {
	c.MoveFront = false
	c.MoveLeft = false
	c.MoveRight = false
	c.MoveBack = false
	c.Interact = false
	c.Attack = false
}
func (c *ControlComponent) Superimpose(state *events.PlayerState) {
	c.MoveFront = state.MoveFront
	c.MoveLeft = state.MoveLeft
	c.MoveRight = state.MoveRight
	c.MoveBack = state.MoveBack
	c.Interact = c.Interact || state.Interact
	c.Attack = state.Attack
	c.Direction = state.Direction
}

func (c *ControlComponent) GetInputVector() geometry.Vec3 {
	input := geometry.Vec3{}
	if c.MoveFront {
		input.Z += 1
	}
	if c.MoveBack {
		input.Z -= 1
	}
	if c.MoveRight {
		input.X -= 1
	}
	if c.MoveLeft {
		input.X += 1
	}
	return input.Normalize()
}
