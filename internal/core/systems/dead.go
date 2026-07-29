package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type DeadSystem struct {
	componentsGetter
	publisher
	broadcaster
	eventQueue chan events.Event
	ship       types.Entity
}

func NewDeadSystem(getter componentsGetter, publisher publisher, subscriber subscriber, broadcaster broadcaster) *DeadSystem {
	s := &DeadSystem{
		getter,
		publisher,
		broadcaster,
		make(chan events.Event, 100),
		"",
	}
	s.SetShip()
	subscriber.Subscribe("dead", s.OnEvent)
	return s
}

func (s *DeadSystem) SetShip() {
	ships := s.GetEntitiesByComponent("ship")
	for sh := range ships {
		s.ship = sh
		break
	}
}

func (s *DeadSystem) Update(dt float32) error {
	return s.drainEvents()
}

func (s *DeadSystem) drainEvents() error {
	for {
		select {
		case e := <-s.eventQueue:
			ev := e.(*events.DeadEvent)
			if ev.ID == s.ship {
				s.Publish(events.NewGameOverEvent(true))
			}
			//s.Publish(events.NewDeleteEvent(ev.ID))
		default:
			return nil
		}
	}
}

func (s *DeadSystem) OnEvent(event events.Event) error {
	select {
	case s.eventQueue <- event:

	default:
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}
