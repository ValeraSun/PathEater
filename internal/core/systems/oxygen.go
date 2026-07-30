package systems

import (
	"log"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	loss            = 1.0 // потеря кислорода в вакууме
	suckForce       = 2.5 // сила притяжения к поломке
	maxSuckDistance = 1.5 // максимальная дистанция притяжения
)

type OxygenSystem struct {
	componentsGetter
	publisher
	players    map[types.Entity]*components.PlayerComponent
	rooms      map[types.Entity]*components.RoomComponent
	breakdowns map[types.Entity]*components.BreakdownComponent
	frameCount int // для ограничения логов
}

func NewOxygenSystem(getter componentsGetter, publisher publisher) *OxygenSystem {
	s := &OxygenSystem{
		getter,
		publisher,
		make(map[types.Entity]*components.PlayerComponent),
		make(map[types.Entity]*components.RoomComponent),
		make(map[types.Entity]*components.BreakdownComponent),
		0,
	}
	s.rooms = s.getRooms()
	log.Println("=== OxygenSystem инициализирован ===")
	log.Printf("Загружено комнат: %d", len(s.rooms))
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
	s.frameCount++
	shouldLog := s.frameCount%60 == 0 // логировать раз в секунду (при 60 FPS)

	s.players = s.getPlayers()
	s.setBreakdowns()

	if shouldLog {
		log.Printf("=== OxygenSystem Update (кадр %d) ===", s.frameCount)
		log.Printf("Игроков: %d, Поломок: %d", len(s.players), len(s.breakdowns))
	}

	for id, player := range s.players {
		room, exists := s.rooms[player.RoomID]
		if !exists || room == nil {
			if shouldLog {
				log.Printf("Игрок %v не находится ни в одной комнате (RoomID: %v)", id, player.RoomID)
			}
			continue
		}

		if shouldLog {
			log.Printf("Игрок %v в комнате %s (вакуум: %v, поломка: %v)",
				id, room.Name, room.Vacuum, room.HasBreakdown)
		}

		// 1. Если есть поломка и НЕТ вакуума - притягиваем к поломке
		if room.HasBreakdown && !room.Vacuum {
			if shouldLog {
				log.Printf("Применяем силу притяжения для игрока %v в комнате %s", id, room.Name)
			}
			s.applySuckForce(id, player)
		} else {
			// Если нет поломки или уже вакуум - сбрасываем силу притяжения
			if shouldLog && room.HasBreakdown && room.Vacuum {
				log.Printf("Сбрасываем силу притяжения для игрока %v (вакуум в комнате %s)", id, room.Name)
			}
			s.clearSuckForce(id)
		}

		// 2. Если вакуум - теряем кислород
		if room.Vacuum {
			if shouldLog {
				log.Printf("Вакуум в комнате %s! Игрок %v теряет кислород", room.Name, id)
			}
			s.applyVacuumDamage(id)
		}
	}

	if shouldLog {
		log.Printf("=== Конец OxygenSystem Update ===")
	}

	return nil
}

// applySuckForce применяет силу притяжения к поломке
func (s *OxygenSystem) applySuckForce(playerID types.Entity, player *components.PlayerComponent) {
	// Получаем externalVelocity
	c, ok := s.GetComponent(playerID, "externalVelocity")
	if !ok {
		log.Printf("applySuckForce: игрок %v не имеет externalVelocity", playerID)
		return
	}
	ext, ok := c.(*components.ExternalVelocityComponent)
	if !ok || ext == nil {
		log.Printf("applySuckForce: externalVelocity игрока %v имеет неверный тип", playerID)
		return
	}

	// Получаем transform (позиция игрока)
	c, ok = s.GetComponent(playerID, "transform")
	if !ok {
		log.Printf("applySuckForce: игрок %v не имеет transform", playerID)
		return
	}
	tr, ok := c.(*components.TransformComponent)
	if !ok || tr == nil {
		log.Printf("applySuckForce: transform игрока %v имеет неверный тип", playerID)
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
		log.Printf("applySuckForce: поломка в комнате %v не найдена, сбрасываем силу", player.RoomID)
		ext.Direction = geometry.Vec3{}
		return
	}

	// Вычисляем направление к поломке
	toBreakdown := targetBreakdown.Position.Sub(tr.Position)
	distance := toBreakdown.Length()

	log.Printf("applySuckForce: игрок %v, расстояние до поломки: %.2f", playerID, distance)

	if distance < 0.001 {
		// Игрок уже на месте поломки - сбрасываем силу
		log.Printf("applySuckForce: игрок %v уже на месте поломки, сбрасываем силу", playerID)
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
	force := suckForce * forceMultiplier
	ext.Direction = direction.Scale(force)
	ext.Direction.Y = 0

	log.Printf("applySuckForce: игрок %v, сила: %.2f, множитель: %.2f, расстояние: %.2f",
		playerID, force, forceMultiplier, distance)
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
		log.Printf("applyVacuumDamage: игрок %v не имеет oxygen", playerID)
		return
	}
	ox, ok := c.(*components.OxygenComponent)
	if !ok || ox == nil {
		log.Printf("applyVacuumDamage: oxygen игрока %v имеет неверный тип", playerID)
		return
	}

	oldOxygen := ox.Oxygen
	isGasp := ox.Leak(loss)

	log.Printf("applyVacuumDamage: игрок %v, кислород: %.2f -> %.2f (потеря: %.2f)",
		playerID, oldOxygen, ox.Oxygen, loss)

	if isGasp {
		log.Printf("applyVacuumDamage: игрок %v задыхается! Наносим урон", playerID)
		s.publisher.Publish(events.NewDamageDealEvent(playerID, 0.1))
	}
}

func (s *OxygenSystem) setBreakdowns() {
	oldCount := len(s.breakdowns)
	s.breakdowns = make(map[types.Entity]*components.BreakdownComponent)

	bs := s.GetEntitiesByComponent("breakdown")
	for id, b := range bs {
		if breakdown, ok := b.(*components.BreakdownComponent); ok && breakdown != nil {
			s.breakdowns[id] = breakdown
		}
	}

	if len(s.breakdowns) != oldCount {
		log.Printf("setBreakdowns: поломок было %d, стало %d", oldCount, len(s.breakdowns))
	}
}
