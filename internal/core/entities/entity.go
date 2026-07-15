package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type EntityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
}

func NewPlayer(adder EntityAdder, id string) types.Entity {
	transform := components.NewTransformComponent()
	velocity := components.NewVelocityComponent()
	control := components.NewControlComponent(id)
	entity, _ := adder.AddEntity(transform, velocity, control) //добавить обработку

	return entity
}

type EntityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
}

func NewBox(width, depth, height float64, adder EntityAdder) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(geometry.GetZeroVector(), geometry.GetZeroVector()),
		components.NewColliderComponent(geometry.NewBoxCollider(
			geometry.GetZeroVector(),
			geometry.Vec3{X: width / 2, Y: height / 2, Z: depth / 2})),
	)

	return e
}

func NewPlayer(adder EntityAdder) types.Entity {
	e, _ := adder.AddEntity(
		components.NewControlComponent(),
		components.NewMovementComponent(10),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewColliderComponent(geometry.NewBoxCollider(
			geometry.GetZeroVector(),
			geometry.Vec3{X: 100, Y: 100, Z: 100})),
		components.NewMovableComponent(),
		components.NewTransformComponent(
			geometry.GetZeroVector(),
			geometry.GetZeroVector(),
		),
	)
	return e

}
