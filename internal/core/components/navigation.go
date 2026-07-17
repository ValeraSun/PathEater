package components

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type NavigationShipComponent struct {
	Position  geometry.Vec2
	Direction geometry.Vec2
	Velocity  geometry.Vec2
}

func (*NavigationShipComponent) Type() string {
	return "navigation_ship"
}

func NewNavigationShipComponent(pos, dir, vel geometry.Vec2, rot float64) *NavigationShipComponent {
	return &NavigationShipComponent{
		Position:  pos,
		Direction: dir,
		Velocity:  vel,
	}
}

func (c *ControlComponent) GetInputVector2() geometry.Vec2 {
	input := geometry.Vec2{}
	if c.MoveFront {
		input.X += 1
	}
	if c.MoveBack {
		input.X -= 1
	}
	if c.MoveRight {
		input.Y += 1
	}
	if c.MoveLeft {
		input.Y -= 1
	}
	return input.Normalize()
}

func (c *NavigationShipComponent) ApplyControl(state *ControlComponent) {
	forwardVector := state.GetInputVector2().Normalize()

	c.Direction = geometry.CombineVectors2(c.Direction, forwardVector)
}
