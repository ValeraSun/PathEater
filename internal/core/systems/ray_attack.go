package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type RayAttackSystem struct {
	getter    componentsGetter
	publisher publisher
}

func NewRayAttackComponent(getter componentsGetter, publisher publisher) *RayAttackSystem {
	return &RayAttackSystem{
		getter:    getter,
		publisher: publisher,
	}
}

func (s *RayAttackSystem) Update(dt float32) error {
	hitboxes := s.getter.GetEntitiesByComponent("hitbox")

	attackers := s.getter.GetEntitiesByComponent("rayAttack")

	for attackerID, c := range attackers {
		attacker := c.(*components.RayAttackComponent)
		attacker.ReduseCooldown(dt)

		if !s.getter.HasComponents(attackerID, "control", "interactionDetector", "transform") {
			continue
		}

		c, _ = s.getter.GetComponent(attackerID, "control")
		control := c.(*components.ControlComponent)

		if control.Interact && attacker.CanAttack() {
			c, _ = s.getter.GetComponent(attackerID, "interactionDetector")
			inter := c.(*components.InteractionDetectorComponent)

			c, _ = s.getter.GetComponent(attackerID, "transform")
			t := c.(*components.TransformComponent)

			ray := inter.GetRay(t.Position, t.Direction)

			for id, c := range hitboxes {
				h := c.(*components.HitboxComponent)

				result := h.Collider.Collide(ray)

				if result.HasCollision && id != attackerID {
					attacker.ResetCooldown()
					e := events.NewDamageDealEvent(id, attacker.Damage)
					s.publisher.Publish(e)
				}
			}

		}
	}

	return nil
}
