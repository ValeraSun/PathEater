package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type AISystem struct {
	componentsGetter
	publisher
}

func NewAISystem(getter componentsGetter, publisher publisher) *AISystem {
	return &AISystem{
		getter,
		publisher,
	}
}

func (s *AISystem) Update(dt float32) error {
	aisRaw := s.GetEntitiesByComponent("ai")

	type enemy struct {
		ai *components.AIComponent
		*components.AttackComponent
		*components.VisionComponent
		*components.TransformComponent
		privMtv geometry.Vec3
		id      types.Entity
	}

	enemies := make([]enemy, 0, len(aisRaw))

	for id, c := range aisRaw {

		ai := c.(*components.AIComponent)

		if s.HasComponents(id, "vision", "collider", "transform", "attack") {
			c, _ := s.GetComponent(id, "vision")
			vision := c.(*components.VisionComponent)

			c, _ = s.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)

			c, _ = s.GetComponent(id, "collider")
			collider := c.(*components.ColliderComponent)

			c, _ = s.GetComponent(id, "attack")
			attack := c.(*components.AttackComponent)

			enemies = append(enemies, enemy{
				ai,
				attack,
				vision,
				transform,
				collider.PrivMTV,
				id,
			})
		}

	}

	for _, enemy := range enemies {

		enemy.ReduceCooldown(dt)

		enemy.ai.MoveFront = false
		distant := enemy.WantPossition.Sub(enemy.Position)
		distant.Y = 0

		if enemy.GoingToLastSee {
			if distant.Length() > cameDistant {
				enemy.ai.MoveFront = true
			}
		} else {

			if distant.Length() > walkDistant {
				enemy.ai.MoveFront = true
			}

			if distant.Length() <= attackDistant {
				e, isAttack := enemy.Attack(enemy.id, enemy.Position, distant.Normalize())

				if isAttack {
					anime := events.NewAlienAttackAnimationEvent(enemy.id)
					s.Publish(anime)
					s.Publish(e)
				}
			}
		}
		enemy.ai.Direction = distant.Normalize().Add(enemy.privMtv.Normalize()).Normalize()

	}

	return nil
}
