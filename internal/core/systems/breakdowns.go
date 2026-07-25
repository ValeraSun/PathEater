package systems

import (
	"fmt"
	"math/rand"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type BreakdownSystem struct {
	getter        componentsGetter
	publisher     publisher
	externalWalls []types.Entity
	eventQueue    chan events.Event
}

func NewBreakSystem(getter componentsGetter, publisher publisher, subscriber subscriber, externalWalls []types.Entity) *BreakdownSystem {
	s := &BreakdownSystem{
		getter:        getter,
		publisher:     publisher,
		externalWalls: externalWalls,
		eventQueue:    make(chan events.Event, 100),
	}
	subscriber.Subscribe("breakdown", s.OnEvent)
	return s
}

func (s *BreakdownSystem) Update(dt float32) error {
	for {
		select {
		case e := <-s.eventQueue:
			ev := e.(*events.BreakdownEvent)
			wallID, pos, normal := s.generateBreakdown()
			s.publisher.Publish(events.NewCreateBreakdownEvent(wallID, pos))
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
	idx := rand.Intn(len(s.externalWalls))
	wallEntity := s.externalWalls[idx]

	colliderComp, _ := s.getter.GetComponent(wallEntity, "collider")
	collider := colliderComp.(*components.ColliderComponent)
	center := collider.Collider.GetCenter()
	box, _ := collider.Collider.(*geometry.BoxCollider)
	normal := box.GetInwardNormal(shipCenter)

	return s.externalWalls[idx], center, normal
}

const alienShift = 0.5

func getAlienPos(pos, normal geometry.Vec3) geometry.Vec3 {
	posRaw := pos.Add(normal.Scale(alienShift))
	posRaw.Z = -1
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
