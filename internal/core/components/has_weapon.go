package components

import "github.com/ValeraSun/PathEater/internal/core/types"

type HasWeaponComponent struct {
	types.Entity
}

func (*HasWeaponComponent) Type() string {
	return "hasWeapon"
}

func NewHasWeaponComponent(weapon types.Entity) *HasWeaponComponent {
	return &HasWeaponComponent{weapon}
}
