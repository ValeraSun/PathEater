package components

import "github.com/ValeraSun/PathEater/internal/core/types"

type WallComponent struct {
	Room         types.Entity
	HasBreakdown bool
}

func (*WallComponent) Type() string {
	return "wall"
}

func NewWallComponent() *WallComponent {
	return &WallComponent{}
}
