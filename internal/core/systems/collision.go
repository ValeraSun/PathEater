package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type CollisionSystem struct {
	getter     componentsGetter
	entitiesID []types.Entity
}

func NewCollisionSystem(getter componentsGetter) *CollisionSystem {
	return &CollisionSystem{
		getter:     getter,
		entitiesID: make([]types.Entity, 0, 128),
	}
}
func (s *CollisionSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("collision")

	for id := range comps {
		s.entitiesID = append(s.entitiesID, id)
	}
	for i := 0; i < len(s.entitiesID); i++ {
		for j := i + 1; j < len(s.entitiesID); j++ {

			id1 := s.entitiesID[i]
			id2 := s.entitiesID[j]
			collision1, _ := comps[id1].(*components.ColliderComponent)
			collision2, _ := comps[id2].(*components.ColliderComponent)

			if s.getter.HasComponents(id1, "transform") && s.getter.HasComponents(id2, "transform") {
				c, _ := s.getter.GetComponent(id1, "transform")
				transform1, _ := c.(*components.TransformComponent)
				c, _ = s.getter.GetComponent(id2, "transform")
				transform2, _ := c.(*components.TransformComponent)

				mtv, isColliding := collision1.Collide(collision2)
				if isColliding {
					if s.getter.HasComponents(id1, "movable") {
						transform1.Position.Add(mtv)
					} else {
						transform2.Position.Add(mtv)
					}
				}

			}
		}
	}

	return nil
}
