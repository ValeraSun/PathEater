package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	shipHealth      = 100
	shipSpeed       = 50
	shipBaggage     = 10
	shipWeaponSpeed = 10
	shipAmmo        = 30
	shipDirX        = 1
	shipDirY        = 0
)

func IsShip(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "ship")
}

func NewShip(adder entityAdder) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			geometry.GetZeroVector(),
			geometry.GetZeroVector(),
		),
		components.NewHealthComponent(shipHealth),
		components.NewMovementComponent(shipSpeed, geometry.Vec3{X: shipDirX, Y: shipDirY, Z: 0}),
		components.NewVelocityComponent(),
		components.NewShipComponent(shipBaggage),
		components.NewWeaponComponent(geometry.GetZeroVector(), shipWeaponSpeed, shipAmmo),
		components.NewColliderComponent(geometry.NewTriangleCollider(
			geometry.Vec3{X: 20, Y: 0, Z: 0},
			geometry.Vec3{X: -20, Y: 20, Z: 0},
			geometry.Vec3{X: -20, Y: -20, Z: 0},
		)),
	)
	return e
}
