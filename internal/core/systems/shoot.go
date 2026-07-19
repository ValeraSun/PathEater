package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type ShootSystem struct {
	getter     componentsGetter
	eventQueue chan *events.ShootEvent
}

func NewShootSystem(getter componentsGetter, subscriber subscriber) *ShootSystem {
	s := &ShootSystem{
		getter:     getter,
		eventQueue: make(chan *events.ShootEvent, 100),
	}
	subscriber.Subscribe("shoot", s.OnEvent)
	return s
}

func (s *ShootSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("bullet")

	s.handleShoot(comps)

	return nil
}

func (s *ShootSystem) handleShoot(comps map[types.Entity]types.Component) {
	for {
		select {
		case e := <-s.eventQueue:
			for id := range comps {
				c, _ := s.getter.GetComponent(id, "bullet")
				bul := c.(*components.BulletComponent)

				c, _ = s.getter.GetComponent(id, "transform")
				transform := c.(*components.TransformComponent)

				c, _ = s.getter.GetComponent(id, "movement")
				mov, _ := c.(*components.MovementComponent)

				c, _ = s.getter.GetComponent(id, "externalVelocity")
				ext, _ := c.(*components.ExternalVelocityComponent)

				ships := s.getter.GetEntitiesByComponent("ship")
				var shipId types.Entity
				for id, _ := range ships {
					shipId = id
				}

				c, _ = s.getter.GetComponent(shipId, "transform")
				trShip := c.(*components.TransformComponent)

				c, _ = s.getter.GetComponent(shipId, "ship")
				ship := c.(*components.ShipComponent)

				ext.Direction = geometry.GetZeroVector().Sub(trShip.Direction.Scale(ship.Speed)).Normalize()
				mov.Direction = e.Direction
				mov.Speed = e.Speed

				x := transform.Position.X
				y := transform.Position.Y
				bul.Visible = x >= -displaySize/2 && x <= displaySize/2 && y >= -displaySize/2 && y <= displaySize/2
			}
		default:
			return
		}
	}
}

func (s *ShootSystem) OnEvent(event events.Event) error {
	se, ok := event.(*events.ShootEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- se:

	default:
		fmt.Printf("Переполена очередь %v\n", s)
	}

	return nil
}
