package systems

import "github.com/ValeraSun/PathEater/internal/core/components"

type ColliderMoveSystem struct {
	getter componentsGetter
}

func NewColliderMoveSystem(getter componentsGetter) *ColliderMoveSystem {
	return &ColliderMoveSystem{
		getter: getter,
	}
}

func (s *ColliderMoveSystem) Update(dt float32) error {
	colliders := s.getter.GetEntitiesByComponent("collider")

	for id, col := range colliders {
		c, exist := s.getter.GetComponent(id, "transform")

		if !exist {
			continue
		}

		t := c.(*components.TransformComponent)

		collider := col.(*components.ColliderComponent)

		collider.Collider.ChangeCenter(t.Position)

	}

	hitboxes := s.getter.GetEntitiesByComponent("hitbox")

	for id, hit := range hitboxes {
		c, exist := s.getter.GetComponent(id, "transform")

		if !exist {
			continue
		}

		t := c.(*components.TransformComponent)

		hitbox := hit.(*components.HitboxComponent)

		hitbox.Collider.ChangeCenter(t.Position)
	}

	return nil
}
