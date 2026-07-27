package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type VelocitySystem struct {
	getter     componentsGetter
	eventQueue chan *events.SetPlayerStateEvent
}

func NewVelocitySystem(getter componentsGetter) *VelocitySystem {
	return &VelocitySystem{
		getter: getter,
	}
}

func (s *VelocitySystem) Update(dt float32) error {
	velocities := s.getter.GetEntitiesByComponent("velocity")

	for entity, vel := range velocities {

		velocity, _ := vel.(*components.VelocityComponent)

		if s.getter.HasComponents(entity, "movement") {
			comp, _ := s.getter.GetComponent(entity, "movement")

			movement, _ := comp.(*components.MovementComponent)

			velocity.Movement = movement.Direction.Scale(movement.Speed * float64(dt))
		}

		if s.getter.HasComponents(entity, "externalVelocity") {
			comp, _ := s.getter.GetComponent(entity, "externalVelocity")

			externalV, _ := comp.(*components.ExternalVelocityComponent)

			velocity.External = externalV.Direction.Scale(float64(dt))
		}
	}

	return nil
}
