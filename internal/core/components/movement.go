package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type MovementComponent struct {
	Speed     float64
	Direction geometry.Vec3
}

func (*MovementComponent) Type() string {
	return "movement"
}

func NewMovementComponent(speed float64) *MovementComponent {
	return &MovementComponent{
		Speed: speed,
	}
}

func (c *MovementComponent) ApplyControl(state *ControlComponent) {
	forwardVector := state.GetInputVector().Normalize()
	direction := state.Direction
	direction.Y = 0
	c.Direction = geometry.CombineVectors(direction, forwardVector)

}
