package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type DoorsSystem struct {
	getter     componentsGetter
	eventQueue chan *events.DoorEvent
}

func NewDoorsSystem(getter componentsGetter, subscriber subscriber) *DoorsSystem {
	s := &DoorsSystem{
		getter:     getter,
		eventQueue: make(chan *events.DoorEvent, 100),
	}
	subscriber.Subscribe("door", s.OnEvent)
	return s
}

func (s *DoorsSystem) Update(dt float32) error {
	for {
		select {
		case e := <-s.eventQueue:

			// Получаем компонент door
			doorComp, err := s.getter.GetComponent(e.ID, "door")
			if err != true {
				continue
			}

			door, ok := doorComp.(*components.DoorComponent)
			if !ok || door == nil {
				continue
			}

			// Получаем компонент collider
			colliderComp, err := s.getter.GetComponent(e.ID, "collider")
			if err != true {
				continue
			}
			collider, ok := colliderComp.(*components.ColliderComponent)
			if !ok || collider == nil {
				continue
			}

			// Переключаем состояние двери
			door.IsOpen = !door.IsOpen
			collider.Enable = !collider.Enable

		default:
			return nil
		}
	}
}

func (s *DoorsSystem) OnEvent(event events.Event) error {
	e, _ := event.(*events.DoorEvent)
	select {
	case s.eventQueue <- e:
	default:
		fmt.Printf("Преполена очередь %v\n", *s)
	}
	return nil
}
