package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type WeaponComponent struct {
	Direction geometry.Vec3
	Speed     float64
	Ammo      int
}

func (*WeaponComponent) Type() string {
	return "weapon"
}

func NewWeaponComponent(direction geometry.Vec3, speed float64, ammo int) *WeaponComponent {
	return &WeaponComponent{
		Direction: direction,
		Speed:     speed,
		Ammo:      ammo,
	}
}
