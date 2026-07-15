package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type CreateSystem struct {
	adder       systemAdder
	getter      componentsGetter
	broadcaster Broadcaster
	eventQueue  chan *events.CreatePlayerEvent
}

func NewCreateSystem(adder systemAdder, getter componentsGetter, broadcaster Broadcaster, subscriber subscriber) *CreateSystem {
	s := &CreateSystem{
		adder:      adder,
		getter:     getter,
		eventQueue: make(chan *events.CreatePlayerEvent, 100),
	}
	subscriber.Subscribe("create", s.OnEvent)
	return s
}

func (s *CreateSystem) Update(dt float32) error {

	s.drainEvents()
<<<<<<< HEAD
=======

>>>>>>> business_logic
	return nil
}

func (s *CreateSystem) drainEvents() {
	for {
		select {
		case e := <-s.eventQueue:
			player := entities.NewPlayer(s.adder)
			s.broadcaster.CreateCameraForPlayer(e.ID, player)
		default:
			return
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
