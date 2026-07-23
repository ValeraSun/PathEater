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
	subscriber.Subscribe("damageShip", s.OnEvent)
	subscriber.Subscribe("damageCosmoAlient", s.OnEvent)
	return s
}

func (s *HealthSystem) Update(dt float32) error {
	for {
		select {
		case e := <-s.eventQueue:
			typ := e.Type()

			switch {
			case typ == "damagePlayer":
				ev, _ := e.(*events.DamagePlayerEvent)
				c, _ := s.getter.GetComponent(ev.ID, "health")
				hp := c.(*components.HealthComponent)
				isDead := hp.Damage(ev.Damage)
				if isDead {
					e := events.NewPlayerDeathEvent(ev.ID)
					s.publisher.Publish(e)
				}
			case typ == "damageShip":
				ev, _ := e.(*events.DamageShipEvent)
				c, _ := s.getter.GetComponent(ev.ID, "health")
				hp := c.(*components.HealthComponent)
				isDead := hp.Damage(ev.Damage)
				if isDead {
					e := events.NewGameOverEvent(false)
					s.publisher.Publish(e)
				}
			case typ == "damageCosmoAlien":
				ev, _ := e.(*events.DamageCosmoAlienEvent)
				c, _ := s.getter.GetComponent(ev.ID, "health")
				hp := c.(*components.HealthComponent)
				isDead := hp.Damage(ev.Damage)
				if isDead {
					e := events.NewDeleteCosmoAlienEvent(ev.ID)
					s.publisher.Publish(e)
				}
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
