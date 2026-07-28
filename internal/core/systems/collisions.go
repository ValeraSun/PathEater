package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type CollisionsSystem struct {
	getter     componentsGetter
	publisher  transfer.EventPublisher
	eventQueue chan events.Event
}

func NewCollisionsSystem(getter componentsGetter, publisher transfer.EventPublisher, subscriber subscriber) *CollisionsSystem {
	s := &CollisionsSystem{
		getter:     getter,
		publisher:  publisher,
		eventQueue: make(chan events.Event, 100),
	}
	subscriber.Subscribe("collision", s.OnEvent)
	return s
}

func (s *CollisionsSystem) Update(dt float32) error {
	return s.drainEvents()
}

func (s *CollisionsSystem) drainEvents() error {
	for {
		select {
		case e := <-s.eventQueue:
			ev := e.(*events.CollisionEvent)
			switch {
			case s.getter.HasComponents(ev.ID1, "ship"):
				switch {
				case s.getter.HasComponents(ev.ID2, "asteroid"):
					s.shipAsteroidCollision(ev.ID1, ev.ID2)
				case s.getter.HasComponents(ev.ID2, "cosmoAlien"):
					s.shipCosmoAlienCollision(ev.ID1, ev.ID2)
				}
			case s.getter.HasComponents(ev.ID1, "asteroid"):
				switch {
				case s.getter.HasComponents(ev.ID2, "ship"):
					s.shipAsteroidCollision(ev.ID2, ev.ID1)
				case s.getter.HasComponents(ev.ID2, "asteroid"):
					s.asteroidAsteroidCollision(ev.ID1, ev.ID2)
				case s.getter.HasComponents(ev.ID2, "cosmoAlien"):
					s.asteroidCosmoAlienCollision(ev.ID1, ev.ID2)
				case s.getter.HasComponents(ev.ID2, "bullet"):
					s.asteroidBulletCollision(ev.ID1, ev.ID2)
				}
			case s.getter.HasComponents(ev.ID1, "cosmoAlien"):
				switch {
				case s.getter.HasComponents(ev.ID2, "ship"):
					s.shipCosmoAlienCollision(ev.ID2, ev.ID1)
				case s.getter.HasComponents(ev.ID2, "asteroid"):
					s.asteroidCosmoAlienCollision(ev.ID2, ev.ID1)
				case s.getter.HasComponents(ev.ID2, "bullet"):
					s.cosmoAlienBulletCollision(ev.ID1, ev.ID2)
				}
			case s.getter.HasComponents(ev.ID1, "bullet"):
				switch {
				case s.getter.HasComponents(ev.ID2, "asteroid"):
					s.asteroidBulletCollision(ev.ID2, ev.ID1)
				case s.getter.HasComponents(ev.ID2, "cosmoAlien"):
					s.cosmoAlienBulletCollision(ev.ID2, ev.ID1)
				}
			}
		default:
			return nil
		}
	}
}

func (s *CollisionsSystem) OnEvent(event events.Event) error {
	select {
	case s.eventQueue <- event:

	default:
		fmt.Printf("Преполена очередь %+v\n", *s)
	}

	return nil
}

func (s *CollisionsSystem) shipAsteroidCollision(shipID, asteroidID types.Entity) {
	s.publisher.Publish(events.NewDamageDealEvent(shipID, 20))
	s.publisher.Publish(events.NewDeleteEvent(asteroidID))
	s.publisher.Publish(events.NewBreakdownEvent(false))
}

func (s *CollisionsSystem) shipCosmoAlienCollision(shipID, alienID types.Entity) {
	s.publisher.Publish(events.NewDamageDealEvent(shipID, 10))
	s.publisher.Publish(events.NewBoardingEvent(alienID))
	s.publisher.Publish(events.NewBreakdownEvent(true))
}

func (s *CollisionsSystem) asteroidAsteroidCollision(asteroid1ID, asteroid2ID types.Entity) {
	s.publisher.Publish(events.NewDeleteEvent(asteroid1ID))
	s.publisher.Publish(events.NewDeleteEvent(asteroid2ID))
}

func (s *CollisionsSystem) asteroidCosmoAlienCollision(asteroidID, alienID types.Entity) {
	s.publisher.Publish(events.NewDeleteEvent(asteroidID))
	s.publisher.Publish(events.NewDamageDealEvent(alienID, 5))
}

func (s *CollisionsSystem) asteroidBulletCollision(asteroidID, bulletID types.Entity) {
	s.publisher.Publish(events.NewDeleteEvent(asteroidID))
	s.publisher.Publish(events.NewDeleteEvent(bulletID))
}

func (s *CollisionsSystem) cosmoAlienBulletCollision(alienID, bulletID types.Entity) {
	s.publisher.Publish(events.NewDamageDealEvent(alienID, 5))
	s.publisher.Publish(events.NewDeleteEvent(bulletID))
}
