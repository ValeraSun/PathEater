package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/events"
)

type DeadSystem struct {
	getter     componentsGetter
	deleter    entityRemover
	eventQueue chan *events.DeadEvent
}

func NewDeadSystem(getter componentsGetter, deleter entityRemover, subscriber subscriber) *DeadSystem {
	s := &DeadSystem{
		getter:     getter,
		deleter:    deleter,
		eventQueue: make(chan *events.DeadEvent, 100),
	}
	subscriber.Subscribe("dead", s.OnEvent)
	return s
}

func (s *DeadSystem) Update(dt float32) error {

	s.drainEvents()

	return nil
}

func (s *DeadSystem) drainEvents() {
	for {
		select {
		case e := <-s.eventQueue:
			if !s.getter.HasComponents(e.ID, "ship") {
				s.deleter.RemoveEntity(e.ID)
			}

		default:
			return
		}
	}
}

func (s *DeadSystem) OnEvent(event events.Event) error {
	ps, ok := event.(*events.DeadEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- ps:

	default:
		fmt.Printf("Переполена очередь %v\n", *s)
	}

	return nil
}
