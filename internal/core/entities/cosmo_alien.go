package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	cosmoAlienRadius = 15
	cosmoAlienSpeed  = 90
	cosmoAlienHealth = 20
)

func IsCosmoAlien(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "cosmoAlien")
}

func NewCosmoAlien(adder entityAdder, pos geometry.Vec3) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			pos,
			geometry.GetZeroVector(),
		),
		components.NewMovementComponent(cosmoAlienSpeed, geometry.GetZeroVector()),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewUpdateComponent(),
		components.NewNavigationEntityComponent(),
		components.NewCosmoAlienComponent(),
		components.NewHealthComponent(cosmoAlienHealth),
		components.NewColliderComponent(geometry.NewCircleCollider(
			pos,
			cosmoAlienRadius,
		)),
		components.NewMovableComponent(),
	)
	return e
}
