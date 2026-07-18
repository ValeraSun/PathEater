package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type VelocityComponent struct {
	Movement geometry.Vec3
	External geometry.Vec3
}

func (*VelocityComponent) Type() string {
	return "velocity"
}

func NewVelocityComponent() *VelocityComponent {
	return &VelocityComponent{}
}

func (c *VelocityComponent) GetTotalVelocity() geometry.Vec3 {
	return c.External.Add(c.Movement)
}

type MovableComponent struct{}

func (*MovableComponent) Type() string {
	return "movable"
}

func NewMovableComponent() *MovableComponent {
	return &MovableComponent{}
}
