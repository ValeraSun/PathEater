package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
)

type ExternalVelocitySystem struct {
	getter componentsGetter
}

func NewExternalVelocitySystem(getter componentsGetter) *ExternalVelocitySystem {
	return &ExternalVelocitySystem{
		getter: getter,
	}
}

func (s *ExternalVelocitySystem) Update(dt float32) error {

	ships := s.getter.GetEntitiesByComponent("ship")

	var shipMove *components.MovementComponent
	for id := range ships {
		c, _ := s.getter.GetComponent(id, "movement")
		shipMove = c.(*components.MovementComponent)
	}

	comps := s.getter.GetEntitiesByComponent("externalVelocity")

	for id, comp := range comps {
		ext, _ := comp.(*components.ExternalVelocityComponent)

		if s.getter.HasComponents(id, "asteroid") {
			ext.Direction = shipMove.Direction.Scale(-1)
		}

	}
	return nil
}
