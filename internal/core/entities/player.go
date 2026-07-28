package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	cameraDistant = 3
	cameraHeight  = 1.9
	cameraOffset  = 0.5

	playerSpeed          = 10
	playerMaxHealth      = 100
	playerHalfHeight     = 1
	playerRadius         = 0.7
	playerSpawnX         = 2
	playerSpawnY         = 1
	playerSpawnZ         = 0
	playerAttackCooldown = 1
	playerDamage         = 50
)

func IsPlayer(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "transform", "control")
}

func NewPlayer(adder entityAdder, clientID string) types.Entity {
	adder.AddEntityByID(
		types.Entity(clientID),
		components.NewInteractionDetectorComponent(cameraDistant, cameraHeight, cameraOffset),
		components.NewControlShipComponent(),
		components.NewPlayerComponent(),
		components.NewUpdateComponent(),
		components.NewControlComponent(clientID),
		components.NewMovementComponent(playerSpeed, geometry.GetZeroVector()),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewHealthComponent(playerMaxHealth),
		components.NewColliderComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{},
			geometry.Vec3{Y: 1},
			playerHalfHeight,
			playerRadius,
		)),
		components.NewHitboxComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{},
			geometry.Vec3{Y: 1},
			playerHalfHeight,
			playerRadius,
		)),
		components.NewRayAttackComponent(playerAttackCooldown, playerDamage),
		components.NewMovableComponent(),
		components.NewTransformComponent(
			geometry.Vec3{
				X: playerSpawnX,
				Y: playerSpawnY,
				Z: playerSpawnZ,
			},
			geometry.GetZeroVector(),
		),
	)
	return types.Entity(clientID)
}
