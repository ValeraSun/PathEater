package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type CreateSystem struct {
	adder       entityAdder
	getter      componentsGetter
	broadcaster Broadcaster
	eventQueue  chan *events.CreatePlayerEvent
}

func NewCreateSystem(adder entityAdder, getter componentsGetter, broadcaster Broadcaster, subscriber subscriber) *CreateSystem {
	s := &CreateSystem{
		adder:       adder,
		getter:      getter,
		broadcaster: broadcaster,
		eventQueue:  make(chan *events.CreatePlayerEvent, 100),
	}
	subscriber.Subscribe("createPlayer", s.OnEvent)
	return s
}

func (s *CreateSystem) Update(dt float32) error {

	s.drainEvents()
	return nil
}

func (s *CreateSystem) drainEvents() {
	for {
		select {
		case e := <-s.eventQueue:
			player := entities.NewPlayer(s.adder, e.ID)
			info := ecs.EntityCreateInfo{
				ID: player,
			}
			s.broadcaster.SendEntityCreate("player", info)
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
