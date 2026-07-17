package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type ColliderComponent2 struct {
	collider geometry.Collider2
}

func (*ColliderComponent2) Type() string {
	return "collider2"
}

func NewColliderComponent2(collider geometry.Collider2) *ColliderComponent2 {
	return &ColliderComponent2{
		collider: collider,
	}
}

func (c *ColliderComponent2) Collide(other *ColliderComponent2) (geometry.Vec2, bool) {
	result := c.collider.Collide(other.collider)
	return result.MTV, result.HasCollision
}
