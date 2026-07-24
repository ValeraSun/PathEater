package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type HealthSystem struct {
	getter     componentsGetter
	publisher  publisher
	eventQueue chan events.Event
}

func NewHealthSystem(getter componentsGetter, publisher publisher, subscriber subscriber) *HealthSystem {
	s := &HealthSystem{
		getter:     getter,
		publisher:  publisher,
		eventQueue: make(chan events.Event, 100),
	}
	subscriber.Subscribe("damageDeal", s.OnEvent)
	return s
}

func (s *HealthSystem) Update(dt float32) error {
	for {
		select {
		case e := <-s.eventQueue:
			ev, _ := e.(*events.DamageDealEvent)
			c, ok := s.getter.GetComponent(ev.ID, "health")
			if !ok {
				return nil
			}
			hp := c.(*components.HealthComponent)
			isDead := hp.Damage(ev.Damage)
			if isDead {
				e := events.NewDeadEvent(ev.ID)
				s.publisher.Publish(e)
			}
		default:
			return nil
		}
	}
}

func (s *HealthSystem) OnEvent(event events.Event) error {
	select {
	case s.eventQueue <- event:
	default:
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}
