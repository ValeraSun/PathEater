package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

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
