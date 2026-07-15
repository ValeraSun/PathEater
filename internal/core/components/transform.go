package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type TransformComponent struct {
	Position  geometry.Vec3
	Direction geometry.Vec3
}

func (*TransformComponent) Type() string {
	return "transform"
}

func NewTransformComponent() *TransformComponent {
	return &TransformComponent{
		Position:  geometry.GetZeroVector(),
		Direction: geometry.GetZeroVector(),
	}
}

type MovableComponent struct {
}

func (*MovableComponent) Type() string {
	return "movable"
}

func NewMovableComponent() *MovableComponent {
	return &MovableComponent{}
}
