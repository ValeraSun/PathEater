package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type NavigationDisplaySystem struct {
	getter componentsGetter
}

func NewNavigationDisplaySystem(getter componentsGetter) *NavigationDisplaySystem {
	s := &NavigationDisplaySystem{
		getter: getter,
	}

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

	comps = s.getter.GetEntitiesByComponent("control")
	for id, c := range comps {
		control := c.(*components.ControlComponent)
		switch {
		case control.Interact && ship.AvailableID != id:
			ship.AvailableID = id
			c, _ = s.getter.GetComponent(id, "controlShip")
			controlShip := c.(*components.ControlShipComponent)
			controlShip.IsControling = true

		case control.Interact && ship.AvailableID == id:
			ship.AvailableID = ""
			c, _ = s.getter.GetComponent(id, "controlShip")
			controlShip := c.(*components.ControlShipComponent)
			controlShip.IsControling = false
		}

	}
	return nil
}
