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
	comps := s.getter.GetEntitiesByComponent("collider")

	type entityData struct {
		id        types.Entity
		collider  *components.ColliderComponent
		transform *components.TransformComponent
	}

	entities := make([]entityData, 0, len(comps))

	for id, comp := range comps {
		collider, ok := comp.(*components.ColliderComponent)
		if !ok || collider == nil {
			continue
		}

		ed := entityData{
			id:       id,
			collider: collider,
		}

		if s.getter.HasComponents(id, "transform") {
			c, _ := s.getter.GetComponent(id, "transform")
			if transform, ok := c.(*components.TransformComponent); ok {
				ed.transform = transform

				collider.Collider.ChangeCenter(transform.Position)
			}
		}

		entities = append(entities, ed)
	}

	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			e1, e2 := entities[i], entities[j]

			mtv, isColliding := e1.collider.Collide(e2.collider)
			if !isColliding {
				continue
			}

			mtv.Y = 0
			switch {
			case e1.transform != nil && e2.transform != nil:

				e1.transform.Position = e1.transform.Position.Add(mtv.Scale(0.5))
				e2.transform.Position = e2.transform.Position.Add(mtv.Scale(-0.5))

			case e1.transform != nil:

				e1.transform.Position = e1.transform.Position.Add(mtv)

			case e2.transform != nil:

				e2.transform.Position = e2.transform.Position.Add(mtv.Scale(-1))
			}
		}
	}

	return nil
}
