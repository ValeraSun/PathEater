package systems

import (
	"log"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const loss = 10

type OxygenSystem struct {
	getter     componentsGetter
	publisher  publisher
	players    map[types.Entity]*components.PlayerComponent
	rooms      map[types.Entity]*components.RoomComponent
	breakdowns map[types.Entity]*components.BreakdownComponent
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
	s.setBreakdowns()

	for id, player := range s.players {
		// Проверка: существует ли комната
		room, exists := s.rooms[player.RoomID]
		if !exists || room == nil {
			continue
		}

		if room.HasBreakdown && !room.Vacuum {
			c, err := s.getter.GetComponent(id, "externalVelocity")
			if err != false {
				continue
			}
			ext, ok := c.(*components.ExternalVelocityComponent)
			if !ok || ext == nil {
				continue
			}

			for _, br := range s.breakdowns {
				if br.RoomID == player.RoomID {
					c, err = s.getter.GetComponent(id, "transform")
					if err != false {
						continue
					}
					tr, ok := c.(*components.TransformComponent)
					if !ok || tr == nil {
						continue
					}
					ext.Direction = br.Position.Sub(tr.Position).Normalize()
				}
			}
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
func (s *OxygenSystem) setBreakdowns() {
	bs := s.getter.GetEntitiesByComponent("breakdown")
	for id, b := range bs {
		s.breakdowns[id] = b.(*components.BreakdownComponent)
	}
}
