package systems

import (
	"log"

	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type DeleteSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
	remover     entitiesRemover
	publisher   publisher
	eventQueue  chan events.Event
}

func NewDeleteSystem(getter componentsGetter, broadcaster Broadcaster, remover entitiesRemover, publisher publisher, subscriber subscriber) *DeleteSystem {
	s := &DeleteSystem{
		getter:      getter,
		publisher:   publisher,
		broadcaster: broadcaster,
		remover:     remover,
		eventQueue:  make(chan events.Event, 100),
	}
	subscriber.Subscribe("deletePlayer", s.OnEvent)
	subscriber.Subscribe("deleteShip", s.OnEvent)
	subscriber.Subscribe("deleteAsteroid", s.OnEvent)
	subscriber.Subscribe("deleteCosmoAlien", s.OnEvent)
	subscriber.Subscribe("deleteBullet", s.OnEvent)
	return s
}

func (s *DeleteSystem) Update(dt float32) error {
	s.drainEvents()
	return nil
}

func (s *DeleteSystem) drainEvents() {
	for {
		select {
		case e := <-s.eventQueue:
			s.handle(e)
		default:
			return
		}
	}
}

func (s *DeleteSystem) handle(e events.Event) {
	typ := e.Type()

	switch typ {
	case "deletePlayer":
		ev, ok := e.(*events.DeletePlayerEvent)
		if !ok {
			logBadType(typ, e)
			return
		}
		if err := entities.SendPlayer(ev.ID, s.getter, s.broadcaster.SendEntityDelete); err != nil {
			logSendFail(typ, err)
		}
		s.remover.RemoveEntity(ev.ID)

	case "deleteShip":
		ev, ok := e.(*events.DeleteShipEvent)
		if !ok {
			logBadType(typ, e)
			return
		}
		if err := entities.SendShip(ev.ID, s.getter, s.broadcaster.SendEntityDelete); err != nil {
			logSendFail(typ, err)
		}
		s.remover.RemoveEntity(ev.ID)

	case "deleteAsteroid":
		ev, ok := e.(*events.DeleteAsteroidEvent)
		if !ok {
			logBadType(typ, e)
			return
		}
		if err := entities.SendAsteroid(ev.ID, s.getter, s.broadcaster.SendEntityDelete); err != nil {
			logSendFail(typ, err)
		}
		s.remover.RemoveEntity(ev.ID)

	case "deleteCosmoAlien":
		ev, ok := e.(*events.DeleteCosmoAlienEvent)
		if !ok {
			logBadType(typ, e)
			return
		}
		if !s.getter.HasComponents(ev.ID, "cosmoAlien") {
			return
		}
		if err := entities.SendCosmoAlien(ev.ID, s.getter, s.broadcaster.SendEntityDelete); err != nil {
			logSendFail(typ, err)
		}
		s.remover.RemoveEntity(ev.ID)

		//e := events.NewCreateAlienEvent(geometry.Vec3{X: 3, Y: 2, Z: -1})
		s.publisher.Publish(e)

	case "deleteBullet":
		ev, ok := e.(*events.DeleteBulletEvent)
		if !ok {
			logBadType(typ, e)
			return
		}
		if err := entities.SendBullet(ev.ID, s.getter, s.broadcaster.SendEntityDelete); err != nil {
			logSendFail(typ, err)
		}
		s.remover.RemoveEntity(ev.ID)

	default:
		log.Printf("DeleteSystem: неизвестный тип события %q", typ)
	}
}

func logBadType(typ string, e events.Event) {
	log.Printf("DeleteSystem: событие %q имеет неожиданный тип %T", typ, e)
}

func logSendFail(typ string, err error) {
	//log.Printf("DeleteSystem: не удалось отправить удаление (%s): %v", typ, err)
}

func (s *DeleteSystem) OnEvent(event events.Event) error {
	select {
	case s.eventQueue <- event:
	default:
		//log.Printf("DeleteSystem: переполнена очередь событий, событие %q отброшено", event.Type())
	}
	return nil
}
