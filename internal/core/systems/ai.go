package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type AISystem struct {
	getter componentsGetter
}

func NewAISystem(getter componentsGetter) *AISystem {
	return &AISystem{
		getter: getter,
	}
}

func (s *AISystem) Update(dt float32) error {
	aisRaw := s.getter.GetEntitiesByComponent("ai")

	type enemy struct {
		ai      *components.AIComponent
		vision  *components.VisionComponent
		privMtv geometry.Vec3
	}

	enemies := make([]enemy, 0, len(aisRaw))

	for id, c := range aisRaw {

		ai := c.(*components.AIComponent)

		if s.getter.HasComponents(id, "vision") && s.getter.HasComponents(id, "collider") {
			comp, _ := s.getter.GetComponent(id, "vision")
			vision := comp.(*components.VisionComponent)

			comp, _ = s.getter.GetComponent(id, "collider")
			collider := comp.(*components.ColliderComponent)

			enemies = append(enemies, enemy{
				ai:      ai,
				vision:  vision,
				privMtv: collider.PrivMTV,
			})
		}

	}

	for _, enemy := range enemies {
		if enemy.vision.CanSee {

			enemy.ai.Direction = enemy.vision.Direction.Normalize()
			enemy.ai.MoveFront = true
		} else {
			enemy.ai.MoveFront = false
		}
	}

	return nil
}
