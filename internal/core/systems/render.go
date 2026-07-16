package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
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

type PlayerData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
}

func (s *RenderSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("update")

	for id := range comps {
		if s.getter.HasComponents(id, "transform") && s.getter.HasComponents(id, "player") {
			c, _ := s.getter.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)

			c, _ = s.getter.GetComponent(id, "control")
			control := c.(*components.ControlComponent)
			id = types.Entity(control.ClientID)

			s.broadcaster.SendEntityUpdate(ecs.EntityInfo{
				ID:   id,
				Type: "player",
				Data: PlayerData{
					Position: transform.Position,
					Rotation: transform.Direction,
				},
			},
			)
		}

	}

	return nil
}
