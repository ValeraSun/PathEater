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
					e1, e2 := shipAsteroidCollision(ev.ID1, ev.ID2)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				case s.getter.HasComponents(ev.ID2, "cosmoAlien"):
					e1, e2 := shipCosmoAlienCollision(ev.ID1, ev.ID2)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				}
			case s.getter.HasComponents(ev.ID1, "asteroid"):
				switch {
				case s.getter.HasComponents(ev.ID2, "ship"):
					e1, e2 := shipAsteroidCollision(ev.ID2, ev.ID1)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				case s.getter.HasComponents(ev.ID2, "asteroid"):
					e1, e2 := asteroidAsteroidCollision(ev.ID1, ev.ID2)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				case s.getter.HasComponents(ev.ID2, "cosmoAlien"):
					e1, e2 := asteroidCosmoAlienCollision(ev.ID1, ev.ID2)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				case s.getter.HasComponents(ev.ID2, "bullet"):
					e1, e2 := asteroidBulletCollision(ev.ID1, ev.ID2)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				}
			case s.getter.HasComponents(ev.ID1, "cosmoAlien"):
				switch {
				case s.getter.HasComponents(ev.ID2, "ship"):
					e1, e2 := shipCosmoAlienCollision(ev.ID2, ev.ID1)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				case s.getter.HasComponents(ev.ID2, "asteroid"):
					e1, e2 := asteroidCosmoAlienCollision(ev.ID2, ev.ID1)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				case s.getter.HasComponents(ev.ID2, "bullet"):
					e1, e2 := cosmoAlienBulletCollision(ev.ID1, ev.ID2)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				}
			case s.getter.HasComponents(ev.ID1, "bullet"):
				switch {
				case s.getter.HasComponents(ev.ID2, "asteroid"):
					e1, e2 := asteroidBulletCollision(ev.ID2, ev.ID1)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				case s.getter.HasComponents(ev.ID2, "cosmoAlien"):
					e1, e2 := cosmoAlienBulletCollision(ev.ID2, ev.ID1)
					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
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

func shipAsteroidCollision(shipID, asteroidID types.Entity) (e1 events.Event, e2 events.Event) {
	return events.NewDamageDealEvent(shipID, 20), events.NewDeleteEvent(asteroidID)
}

func shipCosmoAlienCollision(shipID, alienID types.Entity) (e1 events.Event, e2 events.Event) {
	return events.NewDamageDealEvent(shipID, 10), events.NewBoardingEvent(alienID)
}

func asteroidAsteroidCollision(asteroid1ID, asteroid2ID types.Entity) (e1 events.Event, e2 events.Event) {
	return events.NewDeleteEvent(asteroid1ID), events.NewDeleteEvent(asteroid2ID)
}

func asteroidCosmoAlienCollision(asteroidID, alienID types.Entity) (e1 events.Event, e2 events.Event) {
	return events.NewDeleteEvent(asteroidID), events.NewDamageDealEvent(alienID, 5)
}

func asteroidBulletCollision(asteroidID, bulletID types.Entity) (e1 events.Event, e2 events.Event) {
	return events.NewDeleteEvent(asteroidID), events.NewDeleteEvent(bulletID)
}

func cosmoAlienBulletCollision(alienID, bulletID types.Entity) (e1 events.Event, e2 events.Event) {
	return events.NewDamageDealEvent(alienID, 5), events.NewDeleteEvent(bulletID)
}
