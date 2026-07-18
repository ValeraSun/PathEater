package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type ColliderComponent struct {
	Collider geometry.Collider
}

func (*ColliderComponent) Type() string {
	return "collider"
}

func NewColliderComponent(collider geometry.Collider) *ColliderComponent {
	return &ColliderComponent{
		Collider: collider,
	}
}

func (c *ColliderComponent) Collide(other *ColliderComponent) (geometry.Vec3, bool) {
	result := c.Collider.Collide(other.Collider)
	return result.MTV, result.HasCollision
}
