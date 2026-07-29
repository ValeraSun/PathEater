package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func NewBaggageHitbox(adder entityAdder, x, z float64) types.Entity {
	e, _ := adder.AddEntity(
		components.NewHitboxComponent(geometry.NewBoxCollider(
			geometry.Vec3{
				X: x,
				Y: 1,
				Z: z,
			},
			geometry.Vec3{
				X: 2.4,
				Y: 5,
				Z: 0.65,
			},
			geometry.GetStandartAxes(),
		),
		),
		components.NewTargetComponent(),
		components.NewBaggageComponent(),
	)
	return e
}
