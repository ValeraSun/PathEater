package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	weaponSpeed  = 0.5
	weaponAmmo   = 30
	turningSpeed = 10
	cooldown     = 3
	length       = 2
)

func NewWeapon(adder entityAdder, startAngle, endAngle float64) types.Entity {
	e, _ := adder.AddEntity(
		&components.WeaponComponent{
			StartAngle:   startAngle,
			EndAngle:     endAngle,
			Direction:    startAngle,
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
