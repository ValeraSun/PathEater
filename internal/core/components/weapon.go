package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type WeaponComponent struct {
	AvailableID  string
	Direction    geometry.Vec3
	Speed        float64
	Ammo         int
	ShootSuccess bool
}

func (*WeaponComponent) Type() string {
	return "weapon"
}

func NewWeaponComponent(direction geometry.Vec3, speed float64, ammo int) *WeaponComponent {
	return &WeaponComponent{
		AvailableID: "",
		Direction:   direction,
		Speed:       speed,
		Ammo:        ammo,
	}
}
