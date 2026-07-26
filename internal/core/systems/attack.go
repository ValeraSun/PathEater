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

func NewAttackSystem(getter componentsGetter, subscriber subscriber, publisher publisher) *AttackSystem {
	s := &AttackSystem{
		getter:     getter,
		eventQueue: make(chan *events.AttackEvent, 20),
		publisher:  publisher,
	}
	subscriber.Subscribe("attack", s.OnEvent)
	return s

}

func (s *AttackSystem) Update(dt float32) error {
	hitboxRaw := s.getter.GetEntitiesByComponent("hitbox")

	hitboxes := make([]*entitiesHitbox, 0, len(hitboxRaw))

	for id, c := range hitboxRaw {

		if s.getter.HasComponents(id, "transform") {
			h := c.(*components.HitboxComponent)

			hitboxes = append(hitboxes, &entitiesHitbox{
				hitbox: h,
				id:     id,
			})
		}

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
	for _, hitbox := range hitboxes {
		_, isColliding := hitbox.hitbox.Collide(event.Collider)

		if isColliding && event.Attacker != hitbox.id {
			e := events.NewDamageDealEvent(hitbox.id, event.Damage)
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
