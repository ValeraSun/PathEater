package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type AttackSystem struct {
	getter     componentsGetter
	eventQueue chan *events.AttackEvent
}

func NewAttackSystem(getter componentsGetter) *AttackSystem {
	return &AttackSystem{
		getter: getter,
	}
}

func (s *AttackSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("hitbox")
	hitboxes := make([]*components.HitboxComponent, 0, len(comps))

	for _, c := range comps {
		h := c.(*components.HitboxComponent)
		hitboxes = append(hitboxes, h)
	}

	s.drainEvents(hitboxes)
	return nil
}

func (s *AttackSystem) drainEvents(hitboxes []*components.HitboxComponent) {
	for {
		select {
		case e := <-s.eventQueue:
			handleAttack(e, hitboxes)
		default:
			return
		}
	}
}

func handleAttack(event *events.AttackEvent, hitboxes []*components.HitboxComponent) {
	for _, c := range hitboxes {
		mtv, isColliding := c.Collide(event.Collider)

		if isColliding
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
