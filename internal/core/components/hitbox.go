
import "github.com/ValeraSun/PathEater/internal/core/geometry"

type HitboxComponent struct {
	Collider geometry.Collider
}

func (*HitboxComponent) Type() string {
	return "hitbox"
}

func NewHitboxComponent(collider geometry.Collider) *HitboxComponent {
	return &HitboxComponent{
		Collider: collider,
	}
}

func (c *HitboxComponent) Collide(other *HitboxComponent) (geometry.Vec3, bool) {
	result := c.Collider.Collide(other.Collider)
	return result.MTV, result.HasCollision
}
