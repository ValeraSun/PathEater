package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type ShootSystem struct {
	getter componentsGetter
}

func NewShootSystem(getter componentsGetter) *ShootSystem {
	s := &ShootSystem{
		getter: getter,
	}
	return s
}

func (s *ShootSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("bullet")

	ships := s.getter.GetEntitiesByComponent("ship")
	var shipId types.Entity
	for id := range ships {
		shipId = id
	}

	c, _ := s.getter.GetComponent(shipId, "transform")
	trShip := c.(*components.TransformComponent)

	c, _ = s.getter.GetComponent(shipId, "ship")
	ship := c.(*components.ShipComponent)

	for id, comp := range comps {
		bul := comp.(*components.BulletComponent)

		c, _ := s.getter.GetComponent(id, "transform")
		transform := c.(*components.TransformComponent)

		c, _ = s.getter.GetComponent(id, "externalVelocity")
		ext, _ := comp.(*components.ExternalVelocityComponent)

		ext.Direction = geometry.GetZeroVector().Sub(trShip.Direction.Scale(ship.Speed)).Normalize()

		x := transform.Position.X
		y := transform.Position.Y
		bul.Visible = x >= -displaySize/2 && x <= displaySize/2 && y >= -displaySize/2 && y <= displaySize/2
	}
	return nil
}
