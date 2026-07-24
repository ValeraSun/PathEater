package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type MovementSystem struct {
	getter componentsGetter
}

func NewMovementSystem(getter componentsGetter) *MovementSystem {
	return &MovementSystem{
		getter: getter,
	}
}

func (s *MovementSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("movement")

	for id, comp := range comps {
		movement, _ := comp.(*components.MovementComponent)

		if s.getter.HasComponents(id, "control") && !s.getter.HasComponents(id, "ai", "ship") {
			c, _ := s.getter.GetComponent(id, "control")
			control, _ := c.(*components.ControlComponent)
			movement.ApplyControl(control)
		}

		if !s.getter.HasComponents(id, "control") && s.getter.HasComponents(id, "ai") {
			c, _ := s.getter.GetComponent(id, "ai")
			ai, _ := c.(*components.AIComponent)
			movement.ApplyAI(ai)
		}

		if s.getter.HasComponents(id, "ship") {
			c, _ := s.getter.GetComponent(id, "ship")
			ship := c.(*components.ShipComponent)

			entity := ship.AvailableID

			if entity == "" {
				movement.Direction = geometry.Vec3{X: 1}
				return nil
			}
			if !s.getter.HasComponents(entity, "control") {
				return nil
			}
			ctrlComp, _ := s.getter.GetComponent(entity, "control")
			control := ctrlComp.(*components.ControlComponent)

			if !s.getter.HasComponents(entity, "movement") {
				return nil
			}
			movementRaw, _ := s.getter.GetComponent(id, "movement")
			movement := movementRaw.(*components.MovementComponent)

			movement.ApplyShipControl(control)
		}

	}
	return nil
}
