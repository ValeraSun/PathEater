package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type DeleteSystem struct {
	componentsGetter
	entityRemover
	publisher
	broadcaster
	subscriber

	eventQueue chan *events.DeleteEvent
}

func NewDeleteSystem(getter componentsGetter, remover entityRemover, publisher publisher, subscriber subscriber, broadcaster broadcaster) *DeleteSystem {
	s := &DeleteSystem{
		getter,
		remover,
		publisher,
		broadcaster,
		subscriber,

		make(chan *events.DeleteEvent, 100),
	}
	subscriber.Subscribe("delete", s.OnEvent)
	return s
}

func (s *DeleteSystem) Update(dt float32) error {
	return s.drainEvents()
}

func (s *DeleteSystem) drainEvents() error {
	for {
		select {
		case e := <-s.eventQueue:
			err := s.handle(e)
			if err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (s *DeleteSystem) handle(e *events.DeleteEvent) error {

	err := s.Send(e.ID, s.SendEntityDelete)

	if err != nil {
		return err
	}

	s.RemoveEntity(e.ID)

	return nil
}

func (s *DeleteSystem) OnEvent(event events.Event) error {
	e := event.(*events.DeleteEvent)
	select {
	case s.eventQueue <- e:
	default:
		//log.Printf("DeleteSystem: переполнена очередь событий, событие %q отброшено", event.Type())
	}
	return nil
}
