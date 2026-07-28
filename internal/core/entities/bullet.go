package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func IsBullet(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "bullet")
}

const (
	bulletRadius = 4
	bulletSpeed  = 250
)

func NewBullet(adder entityAdder, pos, dir geometry.Vec3) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			pos,
			dir,
		),
		components.NewBulletComponent(),
		components.NewUpdateComponent(),
		components.NewNavigationEntityComponent(),
		components.NewMovementComponent(bulletSpeed, dir),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewColliderComponent(geometry.NewCircleCollider(
			pos,
			bulletRadius,
		)),
	)
	return e
}
