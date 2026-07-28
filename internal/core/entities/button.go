package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func InitWeapon(adder entityAdder) types.Entity {
	weapon := NewWeapon(adder)

	NewButtonWeaponRight(adder, weapon)
	NewButtonWeaponShoot(adder, weapon)
	NewButtonWeaponLeft(adder, weapon)

	return weapon
}

func NewButtonWeaponRight(adder entityAdder, weapon types.Entity) types.Entity {
	return NewButton(adder,
		geometry.Vec3{X: 29.8, Y: 2.7, Z: -17.5},
		geometry.Vec3{X: 0.52, Y: 0.5, Z: 0.3},
		weapon,
		weaponRight,
	)
}

func NewButtonWeaponShoot(adder entityAdder, weapon types.Entity) types.Entity {
	return NewButton(adder,
		geometry.Vec3{X: 26.59, Y: 2.7, Z: -17.5},
		geometry.Vec3{X: 0.52, Y: 0.5, Z: 0.3},
		weapon,
		weaponShoot,
	)
}

func NewButtonWeaponLeft(adder entityAdder, weapon types.Entity) types.Entity {
	return NewButton(adder,
		geometry.Vec3{X: 23.4, Y: 2.7, Z: -17.5},
		geometry.Vec3{X: 0.52, Y: 0.5, Z: 0.3},
		weapon,
		weaponLeft,
	)
}

func NewButton(adder entityAdder, center, halfExtents geometry.Vec3, weapon types.Entity, interactor func(types.Entity, types.Entity) events.Event) types.Entity {
	e, _ := adder.AddEntity(
		components.NewInteractableComponent(center, halfExtents, "always", interactor),
		components.NewButtonComponent(weapon),
	)
	return e

}

func weaponLeft(source, target types.Entity) events.Event {
	return &events.WeaponEvent{
		Left: true,
	}
}

func weaponRight(source, target types.Entity) events.Event {
	return &events.WeaponEvent{
		Right: true,
	}
}

func weaponShoot(source, target types.Entity) events.Event {
	return &events.WeaponEvent{
		Shooting: true,
	}
}
