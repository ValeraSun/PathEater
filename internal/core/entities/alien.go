package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	alienHalfHeight       = 1
	alienMaxHealth        = 100
	alienRadius           = 0.6
	alienAttackDamage     = 0
	alienAttackCooldown   = 2
	alienAttackDistant    = 1
	alienAttackHalfHeight = 1
	alienAttackRadius     = 1
)

func IsAlien(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "ai")
}

func NewAlien(adder entityAdder, position geometry.Vec3) types.Entity {
	alien, _ := adder.AddEntity(
		components.NewTransformComponent(position, geometry.Vec3{}),
		components.NewVisionComponent(),
		components.NewAttackComponent(
			alienAttackDamage,
			alienAttackCooldown,
			alienAttackDistant,
			geometry.NewCapsuleCollider(
				geometry.GetZeroVector(),
				geometry.Vec3{Y: 1},
				alienAttackHalfHeight,
				alienAttackRadius,
			),
		),
		components.NewAnimationComponent("base"),
		components.NewMovementComponent(2, geometry.GetZeroVector()),
		components.NewMovableComponent(),
		components.NewVelocityComponent(),
		components.NewUpdateComponent(),
		components.NewHealthComponent(alienMaxHealth),
		components.NewAIComponent(),
		components.NewColliderComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{},
			geometry.Vec3{Y: 2},
			alienHalfHeight,
			alienRadius,
		)),
		components.NewHitboxComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{},
			geometry.Vec3{Y: 2},
			alienHalfHeight,
			alienRadius,
		)),
	)
	return alien
}
