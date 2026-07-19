package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type VisionComponent struct {
	CanSee    bool
	Direction geometry.Vec3
	LastSeen  geometry.Vec3
}

func NewVisionComponent() *VisionComponent {
	return &VisionComponent{}
}

func (*VisionComponent) Type() string { return "vision" }
