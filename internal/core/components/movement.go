package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

const shipSpeedMultiplier = 1.5

type MovementComponent struct {
	Speed     float64
	Direction geometry.Vec3
}

func (*MovementComponent) Type() string {
	return "movement"
}

func NewMovementComponent(speed float64, dir geometry.Vec3) *MovementComponent {
	return &MovementComponent{
		Speed:     speed,
		Direction: dir,
	}
}

func (c *MovementComponent) ApplyControl(state *ControlComponent) {
	forwardVector := state.GetInputVector()
	direction := state.Direction
	direction.Y = 0
	direction = direction.Normalize()
	c.Direction = geometry.CombineVectors(direction, forwardVector)
}

func (c *MovementComponent) ApplyAI(state *AIComponent) {
	if state.MoveFront {
		dir := state.Direction
		dir.Y = 0
		c.Direction = dir.Normalize()
	} else {
		c.Direction = geometry.Vec3{}
	}

}

func (c *MovementComponent) ApplyShipControl(state *ControlComponent) {
	forwardVector := geometry.Vec3{X: 1}
	if state.MoveFront {
		forwardVector.Y -= 1
	}
	if state.MoveRight {
		forwardVector.X = forwardVector.X * shipSpeedMultiplier
	}
	if state.MoveBack {
		forwardVector.Y += 1
	}
	c.Direction = forwardVector
}
