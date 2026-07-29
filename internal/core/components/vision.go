package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type Target int

const (
	Nothing Target = iota
	Player
	Baggage
)

type VisionComponent struct {
	CanSee   bool
	Distant  geometry.Vec3
	LastSeen geometry.Vec3
	Target   Target
}

func NewVisionComponent() *VisionComponent {
	return &VisionComponent{}
}

func (*VisionComponent) Type() string { return "vision" }
