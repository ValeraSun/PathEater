package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func IsAsteroid(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "asteroid")
}

func NewAsteroid(adder entityAdder, pos geometry.Vec3, radius float64, dir geometry.Vec3, speed float64) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			pos,
			dir,
		),
		components.NewUpdateComponent(),
		components.NewNavigationEntityComponent(),
		components.NewMovementComponent(speed, dir),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewAsteroidComponent(),
		components.NewColliderComponent(geometry.NewCircleCollider(
			pos,
			radius,
		)),
	)
	return e
}
