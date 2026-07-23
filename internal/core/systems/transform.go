package systems

import "github.com/ValeraSun/PathEater/internal/core/components"

type TransformSystem struct {
	getter componentsGetter
}

func NewTransformSystem(getter componentsGetter) *TransformSystem {
	return &TransformSystem{
		getter: getter,
	}
}

func (s *TransformSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("transform")

	for id, comp := range comps {
		transform, _ := comp.(*components.TransformComponent)

		if s.getter.HasComponents(id, "velocity") {
			c, _ := s.getter.GetComponent(id, "velocity")
			velocity, _ := c.(*components.VelocityComponent)
			transform.Position = transform.Position.Add(velocity.GetTotalVelocity())
		}

		if s.getter.HasComponents(id, "control") {
			c, _ := s.getter.GetComponent(id, "control")
			control, _ := c.(*components.ControlComponent)
			transform.Direction = control.Direction
		}
	}

	return nil
}
