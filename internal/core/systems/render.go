package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
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

type playerData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
	Health   int
}

type shipData struct {
	BaggageStatus int
	Health        int
}

func (s *RenderSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("update")

	for id := range comps {
		if s.getter.HasComponents(id, "transform") && s.getter.HasComponents(id, "control") && s.getter.HasComponents(id, "health") {
			c, _ := s.getter.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)

			c, _ = s.getter.GetComponent(id, "health")
			hp := c.(*components.HealthComponent)

			s.broadcaster.SendEntityUpdate(ecs.EntityInfo{
				ID:   id,
				Type: "player",
				Data: playerData{
					Position: transform.Position,
					Rotation: transform.Direction,
					Health:   hp.Health,
				}})
		}
		if s.getter.HasComponents(id, "ship") {
			c, _ := s.getter.GetComponent(id, "ship")
			ship := c.(*components.ShipComponent)

			c, _ = s.getter.GetComponent(id, "health")
			hp := c.(*components.HealthComponent)

			s.broadcaster.SendEntityUpdate(ecs.EntityInfo{
				ID:   id,
				Type: "ship",
				Data: shipData{
					BaggageStatus: ship.StatusBag,
					Health:        hp.Health,
				}})

		}
	}

	return nil
}
