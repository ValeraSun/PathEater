package components

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type NavigationShipComponent struct {
	position geometry.Vec2
	velocity geometry.Vec2
}

func (*NavigationShipComponent) Type() string {
	return "navigation_ship"
}

func NewNavigationShipComponent(pos, vel geometry.Vec2) *NavigationShipComponent {
	return &NavigationShipComponent{
		position: pos,
		velocity: vel,
	}
}
