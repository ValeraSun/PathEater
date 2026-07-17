package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
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
	comps := s.getter.GetEntitiesByComponent("mesh")

	for id, comp := range comps {
		if s.getter.HasComponents(id, "transform") {

			mesh := comp.(*components.MeshComponent)
			c, _ := s.getter.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)

			s.broadcaster.SendEntityUpdate(ecs.EntityUpdateInfo{
				ID:        id,
				Mesh:      mesh.Path,
				Position:  transform.Position,
				Direction: transform.Direction,
			})
		}
	}	
	
	var navShip *components.NavigationShipComponent
	for id, comp := range compsShip {
		navShip = comp.(*components.NavigationShipComponent)
	}

	comps := s.getter.GetEntitiesByComponent("asteroid")
	for id, comp := range comps {
		asteroid := comp.(*components.AsteroidComponent)
		if s.getter.HasComponents(id, "asteroid") && asteroid.Visible(navShip.Position){
			s.broadcaster.SendEntityUpdate(ecs.EntityUpdateInfo{
				ID:        id,
				//Mesh:      mesh.Path,
				Position:  asteroid.Position.Vec2ToVec3(),
				Direction: asteroid.Direction.Vec2ToVec3(),
			})
		}
	}

	return nil
}
