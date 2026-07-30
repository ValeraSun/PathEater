package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	loss            = 0.0  // кислород не тратится для отладки
	suckForce       = 5.0  // сила притяжения к поломке
	maxSuckDistance = 10.0 // максимальная дистанция притяжения
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

		// 1. Если есть поломка и НЕТ вакуума - притягиваем к поломке
		if room.HasBreakdown && !room.Vacuum {
			s.applySuckForce(id, player)
		} else {
			// Если нет поломки или уже вакуум - сбрасываем силу притяжения
			s.clearSuckForce(id)
		}

		// 2. Если вакуум - теряем кислород
		if room.Vacuum {
			s.applyVacuumDamage(id)
		}
	}
	return nil
}

// applySuckForce применяет силу притяжения к поломке
func (s *OxygenSystem) applySuckForce(playerID types.Entity, player *components.PlayerComponent) {
	// Получаем externalVelocity
	c, ok := s.GetComponent(playerID, "externalVelocity")
	if !ok {
		return
	}
	ext, ok := c.(*components.ExternalVelocityComponent)
	if !ok || ext == nil {
		return
	}

	// Получаем transform (позиция игрока)
	c, ok = s.GetComponent(playerID, "transform")
	if !ok {
		return
	}
	tr, ok := c.(*components.TransformComponent)
	if !ok || tr == nil {
		return
	}

	// Ищем поломку в комнате игрока
	var targetBreakdown *components.BreakdownComponent
	for _, br := range s.breakdowns {
		if br.RoomID == player.RoomID {
			targetBreakdown = br
			break
		}
	}

	if targetBreakdown == nil {
		// Поломка исчезла - сбрасываем силу
		ext.Direction = geometry.Vec3{}
		return
	}

	// Вычисляем направление к поломке
	toBreakdown := targetBreakdown.Position.Sub(tr.Position)
	distance := toBreakdown.Length()

	if distance < 0.001 {
		// Игрок уже на месте поломки - сбрасываем силу
		ext.Direction = geometry.Vec3{}
		return
	}

	// Нормализуем направление
	direction := toBreakdown.Normalize()

	// Сила притяжения зависит от расстояния (чем ближе, тем сильнее)
	forceMultiplier := 1.0
	if distance < maxSuckDistance {
		// Чем ближе к поломке, тем сильнее притяжение
		forceMultiplier = 1.0 + (maxSuckDistance-distance)/maxSuckDistance
	}

	// Ограничиваем максимальную силу
	if forceMultiplier > 3.0 {
		forceMultiplier = 3.0
	}

	// Устанавливаем направление и силу в одном векторе
	ext.Direction = direction.Scale(suckForce * forceMultiplier)
}

// clearSuckForce сбрасывает силу притяжения
func (s *OxygenSystem) clearSuckForce(playerID types.Entity) {
	c, ok := s.GetComponent(playerID, "externalVelocity")
	if !ok {
		return
	}
	ext, ok := c.(*components.ExternalVelocityComponent)
	if !ok || ext == nil {
		return
	}

	// Сбрасываем вектор силы
	ext.Direction = geometry.Vec3{}
}

// applyVacuumDamage применяет урон от вакуума
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
		s.publisher.Publish(events.NewDamageDealEvent(playerID, 1))
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
