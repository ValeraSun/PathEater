package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func NewBreakdown(adder entityAdder, pos geometry.Vec3, wallID types.Entity) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			pos,
			geometry.GetZeroVector(),
		),
		components.NewBreakdownComponent(pos, wallID),
	)
	return e
}
