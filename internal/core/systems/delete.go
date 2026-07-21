package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type DeleteSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
	remover     entitiesRemover
	eventQueue  chan events.Event
}

func NewDeleteSystem(getter componentsGetter, broadcaster Broadcaster, remover entitiesRemover, subscriber subscriber) *DeleteSystem {
	s := &DeleteSystem{
		remover:    remover,
		eventQueue: make(chan events.Event, 100),
	}
	subscriber.Subscribe("deletePlayer", s.OnEvent)
	subscriber.Subscribe("deleteShip", s.OnEvent)
	subscriber.Subscribe("deleteAsteroid", s.OnEvent)
	subscriber.Subscribe("deleteCosmoAlien", s.OnEvent)
	subscriber.Subscribe("deleteBullet", s.OnEvent)
	return s
}

func (s *DeleteSystem) Update(dt float32) error {
	return s.drainEvents()
}

func (s *DeleteSystem) drainEvents() error {
	for {
		var err error

		select {
		case e := <-s.eventQueue:
			typ := e.Type()

			switch {
			case typ == "deletePlayer":
				ev, _ := e.(*events.DeletePlayerEvent)
				s.remover.RemoveEntity(ev.ID)
				err = entities.SendPlayer(ev.ID, s.getter, s.broadcaster.SendEntityDelete)
			case typ == "deleteShip":
				ev, _ := e.(*events.DeleteShipEvent)
				s.remover.RemoveEntity(ev.ID)
				err = entities.SendShip(ev.ID, s.getter, s.broadcaster.SendEntityDelete)
			case typ == "deleteAsteroid":
				ev, _ := e.(*events.DeleteAsteroidEvent)
				s.remover.RemoveEntity(ev.ID)
				err = entities.SendAsteroid(ev.ID, s.getter, s.broadcaster.SendEntityDelete)
			case typ == "deleteCosmoAlien":
				ev, _ := e.(*events.DeleteCosmoAlienEvent)
				s.remover.RemoveEntity(ev.ID)
				err = entities.SendCosmoAlien(ev.ID, s.getter, s.broadcaster.SendEntityDelete)
			case typ == "deleteBullet":
				ev, _ := e.(*events.DeleteBulletEvent)
				s.remover.RemoveEntity(ev.ID)
				err = entities.SendBullet(ev.ID, s.getter, s.broadcaster.SendEntityDelete)
			}
		default:
			return nil
		}

		if err != nil {
			return err
		}
	}
}

func (s *DeleteSystem) OnEvent(event events.Event) error {
	select {
	case s.eventQueue <- event:

	default:
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}
