package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func NewNavigationShip(adder entityAdder) types.Entity {
	e := adder.AddEntity(
		components.NewNavigationShipComponent(
			geometry.GetZeroVector(),
			geometry.GetZeroVector(),
			geometry.GetZeroVector(),
		),
	)
	return e

}

func NewAsteroid(adder entityAdder) types.Entity {
	e := adder.AddEntity(
		components.NewTransformComponent(
			geometry.GetZeroVector(),
			geometry.GetZeroVector(),
		),
		),
	return e

}
