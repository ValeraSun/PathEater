package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type NavigationSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
}

func NewNavigationSystem(getter componentsGetter, broadcaster Broadcaster) *NavigationSystem {
	return &NavigationSystem{
		getter:      getter,
		broadcaster: broadcaster,
	}
}

func (s *NavigationSystem) Update(dt float32) error {
	compsShip := s.getter.GetEntitiesByComponent("navigation_ship")

	var navShip *components.NavigationShipComponent
	var idShip types.Entity
	for id, comp := range compsShip {
		idShip = id
		navShip = comp.(*components.NavigationShipComponent)
	}

	c, _ := s.getter.GetComponent(idShip, "control")
	control, _ := c.(*components.ControlComponent)
	navShip.ApplyControl(control)

	s.broadcaster.SendEntityUpdate(ecs.EntityUpdateInfo{
		ID:        idShip,
		Position:  navShip.Position.Vec2ToVec3(),
		Direction: navShip.Direction.Vec2ToVec3(),
	})

	compsAsteroid := s.getter.GetEntitiesByComponent("asteroid")

	for id, comp := range compsAsteroid {
		asteroid := comp.(*components.AsteroidComponent)

		asteroid.Direction.RotateAroundPoint(&asteroid.Position, navShip.Position, navShip.Direction.CosOfAngleBetweenVec2(geometry.Vec2{X: 1, Y: 0}))

		s.broadcaster.SendEntityUpdate(ecs.EntityUpdateInfo{
			ID:        id,
			Position:  asteroid.Position.Vec2ToVec3(),
			Direction: asteroid.Direction.Vec2ToVec3(),
		})
	}
	return nil
}
