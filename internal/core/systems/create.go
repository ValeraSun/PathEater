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
	eventQueue  chan events.Event
}

func NewCreateSystem(adder entityAdder, getter componentsGetter, broadcaster Broadcaster, subscriber subscriber) *CreateSystem {
	s := &CreateSystem{
		adder:       adder,
		getter:      getter,
		broadcaster: broadcaster,
		eventQueue:  make(chan events.Event, 100),
	}
	subscriber.Subscribe("createPlayer", s.OnEvent)
	subscriber.Subscribe("createShip", s.OnEvent)
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

			switch {
			case typ == "createPlayer":
				c, _ := e.(*events.CreatePlayerEvent)
				player := entities.NewPlayer(s.adder, c.ID)
				err = entities.SendPlayer(player, s.getter, s.broadcaster)
			case typ == "createShip":
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
	select {
	case s.eventQueue <- event:

	default:
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}
