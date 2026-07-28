package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type BoardingSystem struct {
	componentsGetter
	subscriber
	publisher
	eventQueue chan *events.BoardingEvent
}

func NewBoardingSystem(getter componentsGetter, subscriber subscriber, publisher publisher) *BoardingSystem {
	s := &BoardingSystem{
		getter,
		subscriber,
		publisher,
		make(chan *events.BoardingEvent, 100),
	}

	s.Subscribe("boarding", s.OnEvent)
	return s
}

func (s *BoardingSystem) Update(dt float32) error {
	for {
		select {
		case e := <-s.eventQueue:
			del := events.NewDeleteEvent(e.ID)
			s.Publish(del)

			create := events.NewCreateAlienEvent(geometry.Vec3{X: 3, Y: 2, Z: -1})
			s.Publish(create)
		default:
			return nil
		}
	}
}

func (s *BoardingSystem) OnEvent(event events.Event) error {
	e := event.(*events.BoardingEvent)

	select {
	case s.eventQueue <- e:
	default:
		fmt.Println("Преполена очередь boarding")
	}

	return nil
}
