package components

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type HitboxComponent struct {
	Collider geometry.Collider
	Team     events.Target
}

func (*HitboxComponent) Type() string {
	return "hitbox"
}

func NewHitboxComponent(collider geometry.Collider) *HitboxComponent {
	return &HitboxComponent{
		Collider: collider,
	}
}

func (c *HitboxComponent) Collide(other geometry.Collider) (geometry.Vec3, bool) {
	result := c.Collider.Collide(other)
	return result.MTV, result.HasCollision
}
