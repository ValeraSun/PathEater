package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	breakdownSize   = 1.0
	breakdownHealth = 10
)

func NewBreakdown(adder entityAdder, pos geometry.Vec3, wallID, roomID types.Entity) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			pos,
			geometry.GetZeroVector(),
		),
		components.NewBreakdownComponent(pos, wallID, roomID),
		components.NewHealthComponent(breakdownHealth),
		components.NewHitboxComponent(geometry.NewBoxCollider(
			pos,
			geometry.Vec3{X: breakdownSize, Y: breakdownSize, Z: breakdownSize},
			geometry.GetStandartAxes(),
		)),
	)
	return e
}
