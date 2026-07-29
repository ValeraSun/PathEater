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
}

func NewDeadSystem(getter componentsGetter, publisher publisher, subscriber subscriber, broadcaster broadcaster) *DeadSystem {
	s := &DeadSystem{
		getter,
		publisher,
		broadcaster,
		make(chan events.Event, 100),
	}
	subscriber.Subscribe("dead", s.OnEvent)
	return s
}

func (s *DeadSystem) GetShip() types.Entity {
	ships := s.GetEntitiesByComponent("ship")
	for sh := range ships {
		return sh
	}
	return ""
}

func (s *DeadSystem) Update(dt float32) error {
	return s.drainEvents()
}

func (s *DeadSystem) drainEvents() error {
	for {
		select {
		case e := <-s.eventQueue:
			ev := e.(*events.DeadEvent)
			if ev.ID == s.GetShip() {
				s.Publish(events.NewGameOverEvent(false))
			} else {
				s.Publish(events.NewDeleteEvent(ev.ID))
			}
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
