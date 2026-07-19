package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type StalkerSystem struct {
	getter componentsGetter
}

func NewStalkerSystem(getter componentsGetter) *StalkerSystem {
	return &StalkerSystem{
		getter: getter,
	}
}

// задаёт преследователю направление ПО ПРЯМОЙ к цели
func (s *StalkerSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("stalker")

	for id, comp := range comps {
		stalker, _ := comp.(*components.StalkerComponent)

		if s.getter.HasComponents(id, "transform") && s.getter.HasComponents(stalker.Target, "transform") {
			c, _ := s.getter.GetComponent(id, "transorm")
			transform, _ := c.(*components.TransformComponent)

			c, _ = s.getter.GetComponent(id, "movement")
			mov, _ := c.(*components.MovementComponent)

			c, _ = s.getter.GetComponent(id, "externalVelocity")
			ext, _ := comp.(*components.ExternalVelocityComponent)

			cTarg, _ := s.getter.GetComponent(stalker.Target, "transorm")
			transformTarg, _ := cTarg.(*components.TransformComponent)

			cTarg, _ = s.getter.GetComponent(stalker.Target, "velocity")
			velTarg, _ := cTarg.(*components.VelocityComponent)

			transform.Direction = transformTarg.Position.Sub(transform.Position).Normalize()
			mov.Direction = transformTarg.Position.Sub(transform.Position).Normalize()
			ext.Direction = geometry.GetZeroVector().Sub(velTarg.GetTotalVelocity()).Normalize()
		}
	}

	return nil
}
