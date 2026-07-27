package systems

import (
	"fmt"

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
			ev := e.(*events.DoorEvent)
			d, _ := s.getter.GetComponent(ev.ID, "door")
			door := d.(*components.DoorComponent)
			c, _ := s.getter.GetComponent(ev.ID, "collider")
			col := c.(*components.ColliderComponent)
			door.IsOpen = !door.IsOpen
			col.Enable = !col.Enable
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
