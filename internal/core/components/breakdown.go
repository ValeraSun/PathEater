package components

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type BreakdownComponent struct {
	Position geometry.Vec3
	WallID   types.Entity
}

func (*BreakdownComponent) Type() string {
	return "breakdown"
}

func NewBreakdownComponent(pos geometry.Vec3, wallID types.Entity) *BreakdownComponent {
	return &BreakdownComponent{
		Position: pos,
		WallID:   wallID,
	}
}
