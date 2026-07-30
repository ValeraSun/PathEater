package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type DeadSystem struct {
	componentsGetter
	publisher
	players    map[types.Entity]*components.PlayerComponent
	eventQueue chan events.Event
}

func NewDeadSystem(getter componentsGetter, publisher publisher, subscriber subscriber) *DeadSystem {
	s := &DeadSystem{
		getter,
		publisher,
		make(map[types.Entity]*components.PlayerComponent),
		make(chan events.Event, 100),
	}
	s.setPlayers()
	subscriber.Subscribe("dead", s.OnEvent)
	return s
}

func (s *DeadSystem) setPlayers() {
	players := s.GetEntitiesByComponent("player")
	for id, pl := range players {
		s.players[id] = pl.(*components.PlayerComponent)
	}
}

func (s *DeadSystem) Update(dt float32) error {
	s.setPlayers()
	return s.drainEvents()
}

func (s *DeadSystem) drainEvents() error {
	for {
		select {
		case e := <-s.eventQueue:
			ev := e.(*events.DeadEvent)
			ship, _ := getShip(s.componentsGetter)
			switch {
			case ev.ID == ship:
				s.Publish(events.NewGameOverEvent(false))
			case s.isPlayer(ev.ID):
				s.players[ev.ID].Dead = true
				if s.everyoneIsDead() {
					s.Publish(events.NewGameOverEvent(false))
				}
				s.Publish(events.NewDeleteEvent(ev.ID))
			case s.isBreakdown(ev.ID):
				s.Publish(events.NewDamageDealEvent(ship, -5))
				s.Publish(events.NewDeleteEvent(ev.ID))
			default:
				s.Publish(events.NewDeleteEvent(ev.ID))
			}
		default:
			return nil
		}
	}
}

func (s *DeadSystem) OnEvent(event events.Event) error {
	select {
	case s.eventQueue <- event:

	default:
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}

func (s *DeadSystem) isPlayer(id types.Entity) bool {
	return s.HasComponents(id, "player")
}

func (s *DeadSystem) isBreakdown(id types.Entity) bool {
	return s.HasComponents(id, "breakdown")
}

func (s *DeadSystem) everyoneIsDead() bool {
	for _, pl := range s.players {
		if !pl.Dead {
			return false
		}
	}
	return true
}
