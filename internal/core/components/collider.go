package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type ColliderComponent struct {
	collider geometry.Collider
}

func (*ColliderComponent) Type() string {
	return "collider"
}

func NewColliderComponent(collider geometry.Collider) *ColliderComponent {
	return &ColliderComponent{
		collider: collider,
	}
}

func (c *ColliderComponent) Collide(other *ColliderComponent) (geometry.Vec3, bool) {
	result := c.collider.Collide(other.collider)
	return result.MTV, result.HasCollision
}
