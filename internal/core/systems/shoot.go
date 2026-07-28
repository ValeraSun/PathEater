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

	// ПРОВЕРКА: есть ли корабль
	if shipId == "" {
		return nil
	}

	c, err := s.getter.GetComponent(shipId, "transform")
	if err != true {
		return nil
	}
	trShip, ok := c.(*components.TransformComponent)
	if !ok || trShip == nil {
		return nil
	}

	c, err = s.getter.GetComponent(shipId, "ship")
	if err != true {
		return nil
	}
	ship, ok := c.(*components.ShipComponent)
	if !ok || ship == nil {
		return nil
	}

	for id, comp := range comps {
		bul, ok := comp.(*components.BulletComponent)
		if !ok || bul == nil {
			continue
		}

		c, err := s.getter.GetComponent(id, "transform")
		if err != true {
			continue
		}
		transform, ok := c.(*components.TransformComponent)
		if !ok || transform == nil {
			continue
		}

		c, err = s.getter.GetComponent(id, "externalVelocity")
		if err != true {
			continue
		}
		ext, ok := c.(*components.ExternalVelocityComponent)
		if !ok || ext == nil {
			continue
		}

		ext.Direction = geometry.GetZeroVector().Sub(trShip.Direction.Scale(ship.Speed)).Normalize()

		x := transform.Position.X
		y := transform.Position.Y
		bul.Visible = x >= -displayWidth/2 && x <= displayWidth/2 && y >= -displayHeight/2 && y <= displayHeight/2
	}
	return nil
}
