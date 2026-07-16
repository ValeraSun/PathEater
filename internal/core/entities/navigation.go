package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func NewNavigationShip(adder entityAdder) types.Entity {
	e, _ := adder.AddEntity(
		components.NewNavigationShipComponent(
			geometry.GetZeroVector2(),
			geometry.GetZeroVector2(),
		),
		components.NewColliderComponent2(geometry.NewTriangleCollider(
			geometry.Vec2{X: 100, Y: 100},
			geometry.Vec2{X: 100, Y: 100},
			geometry.Vec2{X: 100, Y: 100}),
		),
	)
	return e
}

func NewAsteroid(adder entityAdder) types.Entity {
	pos, vel, rad := components.RandomAsteroid()
	e, _ := adder.AddEntity(
		components.NewAsteroidComponent(pos, vel),
		components.NewColliderComponent2(geometry.NewCircleCollider(pos, rad)),
	)
	return e
}
