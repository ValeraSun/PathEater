package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type NavigationDisplaySystem struct {
	componentsGetter
	subscriber
	eventQueue chan *events.InteractTerminalEvent
}

func NewNavigationDisplaySystem(getter componentsGetter, subscriber subscriber) *NavigationDisplaySystem {
	s := &NavigationDisplaySystem{
		getter,
		subscriber,
		make(chan *events.InteractTerminalEvent, 100),
	}

	s.subscriber.Subscribe("interactTerminal", s.OnEvent)
	return s
}

func (s *NavigationDisplaySystem) Update(dt float32) error {

	_, ship := getShip(s)

	if ship == nil {
		return nil
	}

	for {
		select {
		case e := <-s.eventQueue:
			if e.ID == ship.AvailableID {
				c, _ := s.GetComponent(e.ID, "controlShip")
				controlShip := c.(*components.ControlShipComponent)
				controlShip.IsControling = false

				ship.AvailableID = ""

			} else {
				if ship.AvailableID != "" {
					c, _ := s.GetComponent(ship.AvailableID, "controlShip")
					controlShip := c.(*components.ControlShipComponent)
					controlShip.IsControling = false
				}

				c, _ := s.GetComponent(e.ID, "controlShip")
				controlShip := c.(*components.ControlShipComponent)
				controlShip.IsControling = true

				ship.AvailableID = e.ID
			}
		default:
			return nil
		}
	}
}

func (s *NavigationDisplaySystem) OnEvent(event events.Event) error {
	e := event.(*events.InteractTerminalEvent)

	select {
	case s.eventQueue <- e:
	default:
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}
