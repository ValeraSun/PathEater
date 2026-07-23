package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type DeathSystem struct {
	getter       componentsGetter
	publisher    publisher
	playersAlive map[types.Entity]bool
}

func NewDeathSystem(getter componentsGetter, publisher publisher, subscriber subscriber) *DeathSystem {
	s := &DeathSystem{
		getter:       getter,
		publisher:    publisher,
		playersAlive: make(map[types.Entity]bool),
	}
	players := s.getter.GetEntitiesByComponent("player")
	for id := range players {
		s.playersAlive[id] = true
	}
	subscriber.Subscribe("playerDeath", s.OnEvent)
	return s
}

func (s *DeathSystem) Update(dt float32) error {
	return nil
}

func (s *DeathSystem) OnEvent(event events.Event) error {
	typ := event.Type()
	switch typ {
	case "playerDeath":
		e := event.(*events.PlayerDeathEvent)
		c, _ := s.getter.GetComponent(e.ID, "player")
		pl := c.(*components.PlayerComponent)
		pl.Died = true
		s.playersAlive[e.ID] = false

		var everyoneDied bool = true
		for _, playerAlive := range s.playersAlive {
			if playerAlive {
				everyoneDied = false
				break
			}
		}
		if everyoneDied {
			e := events.NewGameOverEvent(false)
			s.publisher.Publish(e)
		}
	}

	return nil
}
