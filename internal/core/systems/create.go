package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type CreateSystem struct {
	entityAdder
	componentsGetter
	broadcaster
	eventQueue chan events.Event
}

func NewCreateSystem(adder entityAdder, getter componentsGetter, subscriber subscriber, broadcaster broadcaster) *CreateSystem {
	s := &CreateSystem{
		adder,
		getter,
		broadcaster,
		make(chan events.Event, 100),
	}
	subscriber.Subscribe("createPlayer", s.OnEvent)
	subscriber.Subscribe("createShip", s.OnEvent)
	subscriber.Subscribe("createAlien", s.OnEvent)
	subscriber.Subscribe("createAsteroid", s.OnEvent)
	subscriber.Subscribe("createCosmoAlien", s.OnEvent)
	subscriber.Subscribe("createBullet", s.OnEvent)
	subscriber.Subscribe("createBreakdown", s.OnEvent)
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
			var entity types.Entity

			switch typ {

			case "createPlayer":
				c, _ := e.(*events.CreatePlayerEvent)
				entity = entities.NewPlayer(s, c.ID)
			case "createAlien":
				c, _ := e.(*events.CreateAlienEvent)
				entity = entities.NewAlien(s, c.Position)
			case "createShip":
				entity = entities.NewShip(s)
			case "createAsteroid":
				event := e.(*events.CreateAsteroidEvent)
				entity = entities.NewAsteroid(s, event.Position, event.Radius, event.Direction, event.Speed)
			case "createCosmoAlien":
				event := e.(*events.CreateCosmoAlienEvent)
				entity = entities.NewCosmoAlien(s, event.Position)
			case "createBullet":
				event := e.(*events.CreateBulletEvent)
				entity = entities.NewBullet(s, event.Position, event.Direction)
			case "createBreakdown":
				event := e.(*events.CreateBreakdownEvent)
				entity = entities.NewBreakdown(s, event.Position, event.WallID)
			}

			err := s.Send(entity, s.SendEntityCreate)

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
