package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type NavigationDisplaySystem struct {
	getter     componentsGetter
	subscriber subscriber
	eventQueue chan *events.InteractTerminalEvent
}

func NewNavigationDisplaySystem(getter componentsGetter, subscriber subscriber) *NavigationDisplaySystem {
	s := &NavigationDisplaySystem{
		getter:     getter,
		subscriber: subscriber,
		eventQueue: make(chan *events.InteractTerminalEvent, 100),
	}

	s.subscriber.Subscribe("interactTerminal", s.OnEvent)
	return s
}

func (s *NavigationDisplaySystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("ship")

	if len(comps) == 0 {
		return nil
	}
	var c types.Component
	for _, sh := range comps {
		c = sh
	}

	ship := c.(*components.ShipComponent)

	for {
		select {
		case e := <-s.eventQueue:
			if e.ID == ship.AvailableID {
				c, _ = s.getter.GetComponent(e.ID, "controlShip")
				controlShip := c.(*components.ControlShipComponent)
				controlShip.IsControling = false

				ship.AvailableID = ""

			} else {
				if ship.AvailableID != "" {
					c, _ = s.getter.GetComponent(ship.AvailableID, "controlShip")
					controlShip := c.(*components.ControlShipComponent)
					controlShip.IsControling = false
				}

				c, _ = s.getter.GetComponent(e.ID, "controlShip")
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
