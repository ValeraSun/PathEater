package systems

import (
	"fmt"

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
	return s.drainEvents()
}

func (s *CreateSystem) drainEvents() error {
	for {
		select {
		case e := <-s.eventQueue:
			typ := e.Type()
			var err error

			switch typ {
			case "player":
				player := entities.NewPlayer(s.adder, e.ID)
				err = entities.SendPlayer(player, s.getter, s.broadcaster)
			case "ship":
				ship := entities.NewShip(s.adder)
				err = entities.SendShip(ship, s.getter, s.broadcaster)
			}

			if err != nil {
				return err
			}
		default:
			return nil
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
