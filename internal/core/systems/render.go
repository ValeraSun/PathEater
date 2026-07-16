package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type RenderSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
}

func NewRenderSystem(getter componentsGetter, broadcaster Broadcaster) *RenderSystem {
	return &RenderSystem{
		getter:      getter,
		broadcaster: broadcaster,
	}
}
func (s *RenderSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("update")

	for id := range comps {
		if s.getter.HasComponents(id, "transform") {
			c, _ := s.getter.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)
			if s.getter.HasComponents(id, "control") {
				c, _ := s.getter.GetComponent(id, "control")
				control := c.(*components.ControlComponent)
				id = types.Entity(control.ClientID)
			}

			s.broadcaster.SendEntityUpdate(ecs.EntityUpdateInfo{
				ID:        id,
				Position:  transform.Position,
				Direction: transform.Direction,
			})
		}

	}

	return nil
}
