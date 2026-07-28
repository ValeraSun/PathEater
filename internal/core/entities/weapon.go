package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	weaponSpeed  = 0.5
	weaponAmmo   = 30
	turningSpeed = 360
	cooldown     = 0
	length       = 2
)

func NewWeapon(adder entityAdder) types.Entity {
	e, _ := adder.AddEntity(
		&components.WeaponComponent{
			Direction:    0,
			Speed:        weaponSpeed,
			Ammo:         weaponAmmo,
			Cooldown:     cooldown,
			TurningSpeed: turningSpeed,
			Length:       length,
		},
		components.NewUpdateComponent(),
	)
	return e

}
