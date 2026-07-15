package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
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
		if s.getter.HasComponents(id, "control") && !s.getter.HasComponents(id, "ai") {
			c, _ := s.getter.GetComponent(id, "control")
			control, _ := c.(*components.ControlComponent)
			movement.ApplyControl(control)
		}
	}
	return nil
}
