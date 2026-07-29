package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type AISystem struct {
	getter    componentsGetter
	publisher publisher
}

func NewAISystem(getter componentsGetter, publisher publisher) *AISystem {
	return &AISystem{
		getter:    getter,
		publisher: publisher,
	}
}

func (s *AISystem) Update(dt float32) error {
	aisRaw := s.getter.GetEntitiesByComponent("ai")

	type enemy struct {
		ai        *components.AIComponent
		attack    *components.AttackComponent
		vision    *components.VisionComponent
		transform *components.TransformComponent
		id        types.Entity
		privMtv   geometry.Vec3
	}

	enemies := make([]enemy, 0, len(aisRaw))

	for id, c := range aisRaw {

		ai := c.(*components.AIComponent)

		if s.getter.HasComponents(id, "vision", "collider", "transform", "attack") {
			c, _ := s.getter.GetComponent(id, "vision")
			vision := c.(*components.VisionComponent)

			c, _ = s.getter.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)

			c, _ = s.getter.GetComponent(id, "collider")
			collider := c.(*components.ColliderComponent)

			c, _ = s.getter.GetComponent(id, "attack")
			attack := c.(*components.AttackComponent)

			enemies = append(enemies, enemy{
				ai:        ai,
				attack:    attack,
				vision:    vision,
				transform: transform,
				privMtv:   collider.PrivMTV,
				id:        id,
			})
		}

	}

	for _, enemy := range enemies {

		enemy.attack.ReduceCooldown(dt)

		enemy.ai.MoveFront = false

		if enemy.vision.CanSee {

			enemy.ai.Direction = enemy.vision.Distant.Normalize().Add(enemy.privMtv.Normalize())

			if enemy.vision.Distant.Length() > minDistant {
				enemy.ai.MoveFront = true
			} else {
				e, isAttack := enemy.attack.TryAttack(enemy.id, enemy.transform.Position, enemy.vision.Distant)

				if isAttack {
					s.publisher.Publish(e)
				}
			}
		}

	}

	return nil
}
