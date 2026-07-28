package systems

import (
	"fmt"
	"log"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type DoorsSystem struct {
	getter     componentsGetter
	eventQueue chan events.Event
}

func NewDoorsSystem(getter componentsGetter, subscriber subscriber) *DoorsSystem {
	s := &DoorsSystem{
		getter:     getter,
		eventQueue: make(chan events.Event, 100),
	}
	subscriber.Subscribe("door", s.OnEvent)
	return s
}

func (s *DoorsSystem) Update(dt float32) error {
	for {
		select {
		case e := <-s.eventQueue:
			ev, ok := e.(*events.DoorEvent)
			log.Println("DOOR")
			if !ok || ev == nil {
				continue
			}
			log.Println("DOOR EVENT")

			// Получаем компонент door
			doorComp, err := s.getter.GetComponent(ev.ID, "door")
			if err != true {
				continue
			}
			log.Println("DOOR COMPONENT")
			door, ok := doorComp.(*components.DoorComponent)
			if !ok || door == nil {
				continue
			}

			// Получаем компонент collider
			colliderComp, err := s.getter.GetComponent(ev.ID, "collider")
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

			log.Println("DOOR OPEN")
		default:
			return nil
		}
	}
}

func (s *DoorsSystem) OnEvent(event events.Event) error {
	select {
	case s.eventQueue <- event:
	default:
		fmt.Printf("Преполена очередь %v\n", *s)
	}
	return nil
}
