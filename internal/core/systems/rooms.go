package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type RoomsSystem struct {
	getter  componentsGetter
	rooms   map[types.Entity]*components.RoomComponent
	players []types.Entity
}

func NewRoomsSystem(getter componentsGetter) *RoomsSystem {
	s := &RoomsSystem{
		getter: getter,
	}
	s.rooms = s.getRooms()
	return s
}

func (s *RoomsSystem) getPlayers() []types.Entity {
	var players []types.Entity
	ps := s.getter.GetEntitiesByComponent("player")
	for p := range ps {
		players = append(players, p)
	}
	return players
}

func (s *RoomsSystem) getRooms() map[types.Entity]*components.RoomComponent {
	rooms := make(map[types.Entity]*components.RoomComponent)
	rs := s.getter.GetEntitiesByComponent("room")
	for id, room := range rs {
		rooms[id] = room.(*components.RoomComponent)
	}
	return rooms
}

func (s *RoomsSystem) getPlayersPos() map[types.Entity]geometry.Vec3 {
	positions := make(map[types.Entity]geometry.Vec3)
	for _, player := range s.players {
		c, _ := s.getter.GetComponent(player, "transform")
		tr := c.(*components.TransformComponent)
		positions[player] = tr.Position
	}
	return positions
}

func (s *RoomsSystem) Update(dt float32) error {
	s.players = s.getPlayers()
	playersPos := s.getPlayersPos()
	for pl, pos := range playersPos {
		for id, room := range s.rooms {
			if s.playerInRoom(pos, room) {
				c, _ := s.getter.GetComponent(pl, "player")
				player := c.(*components.PlayerComponent)
				player.RoomID = id
			}
		}
	}
	return nil
}

func (s *RoomsSystem) playerInRoom(pos geometry.Vec3, room *components.RoomComponent) bool {
	return s.floatInTheRange(pos.X, room.MinX, room.MaxX) && s.floatInTheRange(pos.Z, room.MinY, room.MaxY)
}

func (s *RoomsSystem) floatInTheRange(num, min, max float64) bool {
	return num >= min && num <= max
}
