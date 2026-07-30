package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/events"
)

type BaggageSystem struct {
	componentsGetter
	subscriber
	publisher

	eventQueue chan *events.DamageBaggageEvent
}

func NewBaggageSystem(getter componentsGetter, subscriber subscriber, publisher publisher) *BaggageSystem {
	s := &BaggageSystem{
		getter,
		subscriber,
		publisher,

		make(chan *events.DamageBaggageEvent, eventQueueSize),
	}

	s.Subscribe("damageBaggage", s.OnEvent)
	return s
}

func (s *BaggageSystem) Update(dt float32) error {
	_, ship := getShip(s)
	for {
		select {
		case <-s.eventQueue:
			ship.BaggageStatus -= 1
			if ship.BaggageStatus == 0 {
				// e := events.NewGameOverEvent(false)
				// s.Publish(e)
			}

		default:
			return nil
		}
	}
}

func (s *BaggageSystem) OnEvent(event events.Event) error {
	e := event.(*events.DamageBaggageEvent)

	select {
	case s.eventQueue <- e:
	default:
		fmt.Println("Преполена очередь boarding")
	}

	return nil
}
