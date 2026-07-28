package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func InitWeapon(adder entityAdder) {
	weapon := NewWeapon(adder, -60, 60)
	NewButtonWeaponRight(adder, weapon)
	NewButtonWeaponShoot(adder, weapon)
	NewButtonWeaponLeft(adder, weapon)
}

func NewButtonWeaponRight(adder entityAdder, weapon types.Entity) types.Entity {
	return NewButton(adder,
		geometry.Vec3{X: 30, Y: 2.3, Z: -17.4},
		geometry.Vec3{X: 0.4, Y: 1.1, Z: 0.5},
		weapon,
		weaponRight,
	)
}

func NewButtonWeaponShoot(adder entityAdder, weapon types.Entity) types.Entity {
	return NewButton(adder,
		geometry.Vec3{X: 27, Y: 2.3, Z: -17.4},
		geometry.Vec3{X: 0.4, Y: 1.1, Z: 0.5},
		weapon,
		weaponShoot,
	)
}

func NewButtonWeaponLeft(adder entityAdder, weapon types.Entity) types.Entity {
	return NewButton(adder,
		geometry.Vec3{X: 23.7, Y: 2.3, Z: -17.4},
		geometry.Vec3{X: 0.4, Y: 1.1, Z: 0.5},
		weapon,
		weaponLeft,
	)
}

func NewButton(adder entityAdder, center, halfExtents geometry.Vec3, weapon types.Entity, interactor func(types.Entity, types.Entity) events.Event) types.Entity {
	e, _ := adder.AddEntity(
		components.NewInteractableComponent(center, halfExtents, interactor),
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
