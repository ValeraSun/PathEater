package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type AttackSystem struct {
	getter     componentsGetter
	publisher  publisher
	eventQueue chan *events.AttackEvent
}

type entitiesHitbox struct {
	hitbox *components.HitboxComponent
	id     types.Entity
}

func NewAttackSystem(getter componentsGetter, subscriber subscriber) *AttackSystem {
	s := &AttackSystem{
		getter:     getter,
		eventQueue: make(chan *events.AttackEvent, 20),
	}
	subscriber.Subscribe("attack", s.OnEvent)
	return s

}

func (s *AttackSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("hitbox")

	hitboxes := make([]*entitiesHitbox, 0, len(comps))

	for id, c := range comps {
		h := c.(*components.HitboxComponent)
		hitboxes = append(hitboxes, &entitiesHitbox{
			hitbox: h,
			id:     id,
		})
	}

	s.drainEvents(hitboxes)
	return nil
}

func (s *AttackSystem) drainEvents(hitboxes []*entitiesHitbox) {
	for {
		select {
		case e := <-s.eventQueue:
			s.handleAttack(e, hitboxes)
		default:
			return
		}
	}
}

func (s *AttackSystem) handleAttack(event *events.AttackEvent, hitboxes []*entitiesHitbox) {
	for _, c := range hitboxes {
		_, isColliding := c.hitbox.Collide(event.Collider)

		if isColliding {
			e := events.NewDamageDealEvent(c.id, event.Damage)
			s.publisher.Publish(e)
		}
	}
}

func (s *AttackSystem) OnEvent(event events.Event) error {
	ps, ok := event.(*events.AttackEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- ps:

	default:
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}
