package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type CreateSystem struct {
	getter     componentsGetter
	adder      entityAdder
	eventQueue chan *events.CreatePlayerEvent
}

func NewCreateSystem(getter componentsGetter) *CreateSystem {
	return &CreateSystem{
		getter:     getter,
		eventQueue: make(chan *events.CreatePlayerEvent, 100),
	}
}

func (s *CreateSystem) Update(dt float32) error {

	s.drainEvents()
	return nil
}

func (s *CreateSystem) drainEvents() {

loop:
	for {
		select {
		case e := <-s.eventQueue:
			entities.NewPlayer(s.adder, e.ID)
		default:
			break loop
		}
	}
}

func (s *CreateSystem) OnEvent(event events.Event) error {
	ps, ok := event.(*events.CreatePlayerEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- ps:

	default:
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}
