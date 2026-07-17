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

<<<<<<< HEAD
<<<<<<< Updated upstream
func (c *ColliderComponent) Collide(other *ColliderComponent) (geometry.Vec3, bool) {
=======
func (c *ColliderComponent) Collide(other *ColliderComponent) (geometry.Vec, bool) {
>>>>>>> Stashed changes
=======
>>>>>>> connection
	result := c.collider.Collide(other.collider)
	return result.MTV, result.HasCollision
}
