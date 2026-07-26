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
			comp, _ := s.getter.GetComponent(id, "vision")
			vision := comp.(*components.VisionComponent)

			comp, _ = s.getter.GetComponent(id, "collider")
			collider := comp.(*components.ColliderComponent)

			comp, _ = s.getter.GetComponent(id, "transform")
			transform := comp.(*components.TransformComponent)

			comp, _ = s.getter.GetComponent(id, "attack")
			attack := comp.(*components.AttackComponent)

			enemies = append(enemies, enemy{
				ai:        ai,
				attack:    attack,
				vision:    vision,
				privMtv:   collider.PrivMTV,
				transform: transform,
				id:        id,
			})
		}

	}

	for _, enemy := range enemies {

		enemy.attack.ReduceCooldown(dt)

		if enemy.vision.CanSee {
			enemy.ai.Direction = enemy.vision.Direction.Normalize()
			enemy.ai.MoveFront = true

			e, isAttack := enemy.attack.TryAttack(enemy.id, enemy.transform.Position, enemy.transform.Direction)

			if isAttack {
				s.publisher.Publish(e)
			}

		} else {
			enemy.ai.MoveFront = false
		}
	}

	return nil
}
