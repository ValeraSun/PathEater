package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	//"github.com/ValeraSun/PathEater/internal/core/types"
)

const displaySize = 200

type AsteroidSystem struct {
	getter     componentsGetter
	eventQueue chan *events.MeteoriteZoneEvent
}

func NewAsteroidSystem(getter componentsGetter, subscriber subscriber) *AsteroidSystem {
	s := &AsteroidSystem{
		getter:     getter,
		eventQueue: make(chan *events.MeteoriteZoneEvent, 100),
	}
	subscriber.Subscribe("meteoriteZone", s.OnEvent)
	return s
}

var active bool

func (s *AsteroidSystem) Update(dt float32) error {
	if active {
		comps := s.getter.GetEntitiesByComponent("asteroid")
		for id, comp := range comps {
			aster := comp.(*components.AsteroidComponent)
			c, _ := s.getter.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)

			/*c, _ = s.getter.GetComponent(id, "collider")
			col := c.(*components.ColliderComponent)	*/

			x := transform.Position.X
			y := transform.Position.Y
			//rad := col.Collider.(*geometry.CircleCollider).GetRadius()
			aster.Visible = x >= -displaySize/2 && x <= displaySize/2 && y >= -displaySize/2 && y <= displaySize/2
		}

		//спавн астероидов
	}
	return nil
}

func (s *AsteroidSystem) OnEvent(event events.Event) error {
	mz, ok := event.(*events.MeteoriteZoneEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- mz:
		e := <-s.eventQueue
		active = !e.Active
	default:
		fmt.Printf("Переполена очередь %v\n", s)
	}

	return nil
}
