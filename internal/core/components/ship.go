package components

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type NavigationShipComponent struct {
	position  geometry.Vec2
	direction geometry.Vec2
	velocity  geometry.Vec2
	radius    float32
}

func (*NavigationShipComponent) Type() string {
	return "navigation_ship"
}

func NewNavigationShipComponent(pos, dir, vel Vec2) *NavigationShipComponent {
	return &NavigationShipComponent{
		position:  pos,
		direction: dir,
		velocity:  vel,
	}
}

func (c *NavigationShipComponent) Reset() {
	c.position = GetZeroVector2()
    c.direction = GetZeroVector2()
    c.velocity = GetZeroVector2()
}