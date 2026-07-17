package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

const displaySize = 200

type TransformComponent struct {
	Position  geometry.Vec3
	Direction geometry.Vec3
}

func (*TransformComponent) Type() string {
	return "transform"
}

func NewTransformComponent(position, direction geometry.Vec3) *TransformComponent {
	return &TransformComponent{
		Position:  position,
		Direction: direction,
	}
}

func (trC *TransformComponent) Visible() bool {
	pos := trC.Position
	return pos.X >= -displaySize/2 && pos.X <= displaySize/2 && pos.Y >= -displaySize/2 && pos.Y >= displaySize/2
}
