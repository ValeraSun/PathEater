package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type DeadSystem struct {
	getter     componentsGetter
	deleter    entityRemover
	publisher  publisher
	eventQueue chan *events.DeadEvent
}

func NewDeadSystem(getter componentsGetter, deleter entityRemover, subscriber subscriber, publisher publisher) *DeadSystem {
	s := &DeadSystem{
		getter:     getter,
		deleter:    deleter,
		publisher:  publisher,
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
			switch {
			case entities.IsAsteroid(e.ID, s.getter):
				e := events.NewDeleteAsteroidEvent(e.ID)
				s.publisher.Publish(e)
			case entities.IsCosmoAlien(e.ID, s.getter):
				e := events.NewDeleteCosmoAlienEvent(e.ID)
				if e == nil {
					return
				}
				s.publisher.Publish(e)
			case s.getter.HasComponents(e.ID, "ai"):
				e := events.NewDeleteAlienEvent(e.ID)
				if e == nil {
					return
				}
				s.publisher.Publish(e)
			}

			// if !s.getter.HasComponents(e.ID, "ship") {
			// 	s.deleter.RemoveEntity(e.ID)
			// }

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
		fmt.Printf("2 Переполена очередь %+v\n", *s)
	}

	return nil
}
