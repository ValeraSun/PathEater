package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/events"
	/*"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"*/)

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
	//comps := s.getter.GetEntitiesByComponent("bullet")

	//s.handleShoot(comps)

	return nil
}

/*func (s *ShootSystem) handleShoot(comps map[types.Entity]types.Component) {
	for {
		select {
		case e := <-s.eventQueue:
			for id := range comps {
			c, _ := s.getter.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)

			c, _ = s.getter.GetComponent(id, "externalVelocity")
			ext, _ := comp.(*components.ExternalVelocityComponent)

			ext

			}
		default:
			return
		}
	}
}*/

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
