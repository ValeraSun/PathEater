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
	for id, comp := range compsShip {
		navShip = comp.(*components.NavigationShipComponent)
	}

	c, _ := s.getter.GetComponent(idShip, "control")
	control, _ := c.(*components.ControlComponent)
	navShip.ApplyControl(control)

	navShip.Position.Add(navShip.Velocity)

	compsAsteroid := s.getter.GetEntitiesByComponent("asteroid")

	for id, comp := range compsAsteroid {

		asteroid := comp.(*components.AsteroidComponent)

		asteroid.Direction.RotateAroundPoint(&asteroid.Position, navShip.Position, navShip.Direction.CosOfAngleBetweenVec2(geometry.Vec2{X: 1, Y: 0}))
		asteroid.Position.Add(asteroid.Velocity.Sub(navShip.Velocity))
	}
	return nil
}
