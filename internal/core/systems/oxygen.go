package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const loss = 10

type OxygenSystem struct {
	getter    componentsGetter
	publisher publisher
	players   map[types.Entity]*components.PlayerComponent
	rooms     map[types.Entity]*components.RoomComponent
}

func NewOxygenSystem(getter componentsGetter, publisher publisher) *OxygenSystem {
	s := &OxygenSystem{
		getter:    getter,
		publisher: publisher,
	}
	s.players = s.getPlayers()
	s.rooms = s.getRooms()
	return s
}

func (s *OxygenSystem) getPlayers() map[types.Entity]*components.PlayerComponent {
	players := make(map[types.Entity]*components.PlayerComponent)
	ps := s.getter.GetEntitiesByComponent("player")
	for id, p := range ps {
		players[id] = p.(*components.PlayerComponent)
	}
	return players
}

func (s *OxygenSystem) getRooms() map[types.Entity]*components.RoomComponent {
	rooms := make(map[types.Entity]*components.RoomComponent)
	rs := s.getter.GetEntitiesByComponent("room")
	for id, room := range rs {
		rooms[id] = room.(*components.RoomComponent)
	}
	return rooms
}

func (s *OxygenSystem) Update(dt float32) error {
	for id, player := range s.players {
		if s.rooms[player.RoomID].Vacuum {
			c, _ := s.getter.GetComponent(id, "oxygen")
			ox := c.(*components.OxygenComponent)
			isGasp := ox.Leak(loss)
			if isGasp {
				s.publisher.Publish(events.NewDamageDealEvent(id, 1))
			}
		}
	}
	return nil
}
