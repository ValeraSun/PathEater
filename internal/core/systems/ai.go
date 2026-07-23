package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
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
		vision    *components.VisionComponent
		transform *components.TransformComponent
		privMtv   geometry.Vec3
	}

	enemies := make([]enemy, 0, len(aisRaw))

	for id, c := range aisRaw {

		ai := c.(*components.AIComponent)

		if s.getter.HasComponents(id, "vision") && s.getter.HasComponents(id, "collider") && s.getter.HasComponents(id, "transform") {
			comp, _ := s.getter.GetComponent(id, "vision")
			vision := comp.(*components.VisionComponent)

			comp, _ = s.getter.GetComponent(id, "collider")
			collider := comp.(*components.ColliderComponent)

			comp, _ = s.getter.GetComponent(id, "transform")
			transform := comp.(*components.TransformComponent)

			enemies = append(enemies, enemy{
				ai:        ai,
				vision:    vision,
				privMtv:   collider.PrivMTV,
				transform: transform,
			})
		}

	}

	for _, enemy := range enemies {
		if enemy.ai.AttackCooldown-dt > 0 {
			enemy.ai.AttackCooldown -= dt
		} else {
			enemy.ai.AttackCooldown = 0
		}

		if enemy.vision.CanSee {
			enemy.ai.Direction = enemy.vision.Direction.Normalize()
			enemy.ai.MoveFront = true

			if enemy.ai.Direction.Length() < 2 && enemy.ai.AttackCooldown == 0 {
				enemy.ai.AttackCooldown = 1
				e := events.NewAttackEvent(
					20,
					geometry.NewCapsuleCollider(
						enemy.transform.Position.Add(enemy.ai.Direction),
						geometry.Vec3{Y: 1},
						1, 1,
					),
					events.Friend,
				)
				s.publisher.Publish(e)
			}

		} else {
			enemy.ai.MoveFront = false
		}
	}

	return nil
}
