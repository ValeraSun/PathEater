package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	loss            = 1.0
	suckForce       = 2.5
	maxSuckDistance = 1.5
	maxForce        = 3.0
	oxygenDamage    = 0.1
)

type OxygenSystem struct {
	componentsGetter
	publisher
	players    map[types.Entity]*components.PlayerComponent
	rooms      map[types.Entity]*components.RoomComponent
	breakdowns map[types.Entity]*components.BreakdownComponent
}

func NewOxygenSystem(getter componentsGetter, publisher publisher) *OxygenSystem {
	s := &OxygenSystem{
		getter,
		publisher,
		make(map[types.Entity]*components.PlayerComponent),
		make(map[types.Entity]*components.RoomComponent),
		make(map[types.Entity]*components.BreakdownComponent),
	}
	s.rooms = s.getRooms()
	return s
}

func (s *OxygenSystem) getPlayers() map[types.Entity]*components.PlayerComponent {
	players := make(map[types.Entity]*components.PlayerComponent)
	ps := s.GetEntitiesByComponent("player")
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
	rs := s.GetEntitiesByComponent("room")
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
		room, exists := s.rooms[player.RoomID]
		if !exists || room == nil {
			continue
		}

		if room.HasBreakdown && !room.Vacuum {
			s.applySuckForce(id, player)
		} else {
			s.clearSuckForce(id)
		}

		if room.Vacuum {
			s.applyVacuumDamage(id)
		}
	}

	return nil
}

func (s *OxygenSystem) applySuckForce(playerID types.Entity, player *components.PlayerComponent) {
	c, ok := s.GetComponent(playerID, "externalVelocity")
	if !ok {
		return
	}
	ext, ok := c.(*components.ExternalVelocityComponent)
	if !ok || ext == nil {
		return
	}

	c, ok = s.GetComponent(playerID, "transform")
	if !ok {
		return
	}
	tr, ok := c.(*components.TransformComponent)
	if !ok || tr == nil {
		return
	}

	var targetBreakdown *components.BreakdownComponent
	for _, br := range s.breakdowns {
		if br.RoomID == player.RoomID {
			targetBreakdown = br
			break
		}
	}

	if targetBreakdown == nil {
		ext.Direction = geometry.Vec3{}
		return
	}

	toBreakdown := targetBreakdown.Position.Sub(tr.Position)
	distance := toBreakdown.Length()

	if distance < 0.001 {
		ext.Direction = geometry.Vec3{}
		return
	}

	direction := toBreakdown.Normalize()

	forceMultiplier := 1.0
	if distance < maxSuckDistance {
		forceMultiplier = 1.0 + (maxSuckDistance-distance)/maxSuckDistance
	}

	if forceMultiplier > maxForce {
		forceMultiplier = maxForce
	}

	force := suckForce * forceMultiplier
	ext.Direction = direction.Scale(force)
	ext.Direction.Y = 0
}

func (s *OxygenSystem) clearSuckForce(playerID types.Entity) {
	c, ok := s.GetComponent(playerID, "externalVelocity")
	if !ok {
		return
	}
	ext, ok := c.(*components.ExternalVelocityComponent)
	if !ok || ext == nil {
		return
	}

	ext.Direction = geometry.Vec3{}
}

func (s *OxygenSystem) applyVacuumDamage(playerID types.Entity) {
	c, ok := s.GetComponent(playerID, "oxygen")
	if !ok {
		return
	}
	ox, ok := c.(*components.OxygenComponent)
	if !ok || ox == nil {
		return
	}

	isGasp := ox.Leak(loss)

	if isGasp {
		s.publisher.Publish(events.NewDamageDealEvent(playerID, oxygenDamage))
	}
}

func (s *OxygenSystem) setBreakdowns() {
	s.breakdowns = make(map[types.Entity]*components.BreakdownComponent)

	bs := s.GetEntitiesByComponent("breakdown")
	for id, b := range bs {
		if breakdown, ok := b.(*components.BreakdownComponent); ok && breakdown != nil {
			s.breakdowns[id] = breakdown
		}
	}
}
