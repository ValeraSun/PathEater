package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const maxTransform = 5

type CollisionSystem struct {
	getter     componentsGetter
	publisher  transfer.EventPublisher
	entitiesID []types.Entity
}

func NewCollisionSystem(getter componentsGetter, publisher transfer.EventPublisher) *CollisionSystem {
	return &CollisionSystem{
		getter:     getter,
		publisher:  publisher,
		entitiesID: make([]types.Entity, 0, 128),
	}
}

func (s *CollisionSystem) Update(dt float32) error {
	collidersRaw := s.getter.GetEntitiesByComponent("collider")

	colliders := make([](*components.ColliderComponent), 0, len(collidersRaw))

	for id, collider := range collidersRaw {
		if !s.getter.HasComponents(id, "transform") {
			c, _ := collider.(*components.ColliderComponent)

			colliders = append(colliders, c)
		}
	}

	type movable struct {
		id        types.Entity
		transform *components.TransformComponent
		collider  *components.ColliderComponent
	}

	movables := make([]movable, 0, maxTransform)

	for id, collider := range collidersRaw {
		if s.getter.HasComponents(id, "transform", "transform") {
			c, _ := s.getter.GetComponent(id, "transform")
			t, _ := c.(*components.TransformComponent)

			col, _ := collider.(*components.ColliderComponent)

			col.Collider.ChangeCenter(t.Position) //Сразу меняем центр коллайдера

			movables = append(movables, movable{
				id:        id,
				transform: t,
				collider:  col,
			})
		}
	}

	for i := 0; i < len(movables); i++ {
		for j := i + 1; j < len(movables); j++ {

			mtv, isColliding := movables[i].collider.Collide(movables[j].collider)

			if !isColliding {
				continue
			}

			mtv.Y = 0
			movables[i].transform.Position = movables[i].transform.Position.Add(mtv.Scale(0.5))
			movables[j].transform.Position = movables[j].transform.Position.Add(mtv.Scale(-0.5))

			e := events.NewCollisionEvent(movables[i].id, movables[j].id)
			s.publisher.Publish(e)
		}
	}

	for _, m := range movables {
		m.collider.PrivMTV = geometry.Vec3{}
		for _, c := range colliders {

			mtv, isColliding := m.collider.Collide(c)

			if !isColliding || mtv.IsZero() {
				continue

			}

			mtv.Y = 0

			m.transform.Position = m.transform.Position.Add(mtv)
			m.collider.PrivMTV = mtv.Scale(float64(1 / dt))
		}
	}
	return nil
}
