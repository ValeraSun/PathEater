package systems

import (
	"fmt"
	"math/rand"

	"github.com/ValeraSun/PathEater/internal/config"
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type BreakdownSystem struct {
	componentsGetter
	publisher
	externalWalls []types.Entity
	rooms         map[types.Entity]types.Entity
	eventQueue    chan events.Event
}

func NewBreakdownSystem(getter componentsGetter, publisher publisher, subscriber subscriber, externalWalls []types.Entity) *BreakdownSystem {
	s := &BreakdownSystem{
		getter,
		publisher,
		externalWalls,
		make(map[types.Entity]types.Entity),
		make(chan events.Event, 100),
	}
	s.rooms = config.GetRooms()
	subscriber.Subscribe("breakdown", s.OnEvent)
	return s
}

func (s *BreakdownSystem) Update(dt float32) error {
	for {
		select {
		case e := <-s.eventQueue:
			ev := e.(*events.BreakdownEvent)
			wallID, pos, normal := s.generateBreakdown()

			// ПРОВЕРЯЕМ СУЩЕСТВОВАНИЕ
			room := s.rooms[wallID]
			if room == "" {
				continue // если комнаты нет - пропускаем
			}

			c, _ := s.GetComponent(room, "room")
			roomComp := c.(*components.RoomComponent)
			roomComp.HasBreakdown = true
			s.publisher.Publish(events.NewCreateBreakdownEvent(wallID, room, pos))
			s.publisher.Publish(events.NewVacuumRecalculateEvent())
			if ev.SpawnAlien {
				alienPos := getAlienPos(pos, normal)
				s.publisher.Publish(events.NewCreateAlienEvent(alienPos))
			}
		default:
			return nil
		}
	}
}

var shipCenter geometry.Vec3 = geometry.Vec3{X: 24.67, Y: 2.8, Z: -9.0}

func (s *BreakdownSystem) generateBreakdown() (types.Entity, geometry.Vec3, geometry.Vec3) {
	if len(s.externalWalls) == 0 {
		return "", geometry.Vec3{}, geometry.Vec3{}
	}
	idx := rand.Intn(len(s.externalWalls))
	wallEntity := s.externalWalls[idx]

	if !s.HasComponents(s.externalWalls[idx], "collider") {
		return "", geometry.Vec3{}, geometry.Vec3{}
	}

	colliderComp, _ := s.GetComponent(wallEntity, "collider")
	collider := colliderComp.(*components.ColliderComponent)
	center := collider.Collider.GetCenter()
	box, _ := collider.Collider.(*geometry.BoxCollider)
	normal := box.GetInwardNormal(shipCenter)

	return s.externalWalls[idx], center, normal
}

const alienShift = 0.5

func getAlienPos(pos, normal geometry.Vec3) geometry.Vec3 {
	posRaw := pos.Add(normal.Scale(alienShift))
	posRaw.Y = 2
	return posRaw
}

func (s *BreakdownSystem) OnEvent(event events.Event) error {
	select {
	case s.eventQueue <- event:

	default:
		fmt.Printf("Преполена очередь %v\n", *s)
	}

	return nil
}
