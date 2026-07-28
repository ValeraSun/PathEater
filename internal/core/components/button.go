package components

import "github.com/ValeraSun/PathEater/internal/core/types"

type ButtonComponent struct {
	Weapon  types.Entity
	Pressed bool
}

func (*ButtonComponent) Type() string {
	return "button"
}

func NewButtonComponent(weapon types.Entity) *ButtonComponent {
	return &ButtonComponent{
		Weapon: weapon,
	}
}
