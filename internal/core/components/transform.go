package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type TransformComponent struct {
	Position  geometry.Vector3
	Direction geometry.Vector3
}

func (*TransformComponent) Type() string {
	return "move"
}

func NewTransformComponent() *TransformComponent {
	return &TransformComponent{
		Position:  geometry.GetZeroVector(),
		Direction: geometry.GetZeroVector(),
	}
}
