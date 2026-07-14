package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type VelocityComponent struct {
	Direction geometry.Vector3
}

func (*VelocityComponent) Type() string {
	return "velocity"
}

func NewVelocityComponent() *VelocityComponent {
	return &VelocityComponent{
		Direction: geometry.GetZeroVector(),
	}
}

func (c *VelocityComponent) ApplyControl(state *ControlComponent) {
	forwardVector := state.GetInputVector().Normalize()

	c.Direction = geometry.CombineVectors(c.Direction, forwardVector)

}
