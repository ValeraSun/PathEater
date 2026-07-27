package systems

import (
	"log"

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
	s.rooms = s.getRooms()
	return s
}

func (s *OxygenSystem) getPlayers() map[types.Entity]*components.PlayerComponent {
	players := make(map[types.Entity]*components.PlayerComponent)
	ps := s.getter.GetEntitiesByComponent("player")
	for id, p := range ps {
		player, ok := p.(*components.PlayerComponent)
		if ok && player != nil {
			players[id] = player
		}
	}
	return players
}

func (s *OxygenSystem) getRooms() map[types.Entity]*components.RoomComponent {
	rooms := make(map[types.Entity]*components.RoomComponent)
	rs := s.getter.GetEntitiesByComponent("room")
	for id, room := range rs {
		r, ok := room.(*components.RoomComponent)
		if ok && r != nil {
			rooms[id] = r
		}
	}
	return rooms
}

func (s *OxygenSystem) Update(dt float32) error {
	s.players = s.getPlayers()

	for id, player := range s.players {
		// Проверка: существует ли комната
		room, exists := s.rooms[player.RoomID]
		if !exists || room == nil {
			continue
		}

		// Проверка: вакуум в комнате
		if room.Vacuum {
			c, err := s.getter.GetComponent(id, "oxygen")
			if err != false {
				continue
			}

			ox, ok := c.(*components.OxygenComponent)
			if !ok || ox == nil {
				continue
			}

			isGasp := ox.Leak(loss)
			log.Println("КИСЛОРОД: ", ox.Oxygen)
			if isGasp {
				s.publisher.Publish(events.NewDamageDealEvent(id, 1))
			}
		}
	}
	return nil
}
