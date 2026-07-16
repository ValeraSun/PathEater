package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type entityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
}

// func NewBox(width, depth, height float64, adder entityAdder) types.Entity {
// 	e, _ := adder.AddEntity(
// 		components.NewTransformComponent(geometry.GetZeroVector(), geometry.GetZeroVector()),
// 		components.NewColliderComponent(geometry.NewBoxCollider(
// 			geometry.GetZeroVector(),
// 			geometry.Vec3{X: width / 2, Y: height / 2, Z: depth / 2})),
// 	)

// 	return e
// }

func NewPlayer(adder entityAdder, clientID string) types.Entity {
	e, _ := adder.AddEntity(
		components.NewMeshComponent(""),
		components.NewControlComponent(clientID),
		components.NewMovementComponent(10),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewColliderComponent(geometry.NewBoxCollider(
			geometry.GetZeroVector(),
			geometry.Vec3{X: 100, Y: 100, Z: 100},
			[3]geometry.Vec3{
				{X: 1, Y: 0, Z: 0},
				{X: 0, Y: 1, Z: 0},
				{X: 0, Y: 0, Z: 1},
			},
		)),
		components.NewMovableComponent(),
		components.NewTransformComponent(
			geometry.GetZeroVector(),
			geometry.GetZeroVector(),
		),
	)
	return e

}

func NewShip(adder entityAdder) types.Entity {
	e, _ := adder.AddEntity(
		components.NewShipComponent(),
	)
	return e
}
