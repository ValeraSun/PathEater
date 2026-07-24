package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type ShipSystem struct {
	getter     componentsGetter
	eventQueue chan *events.SetShipStateEvent
}

func NewShipSystem(getter componentsGetter, subscriber subscriber) *ShipSystem {
	s := &ShipSystem{
		getter: getter,
	}
	return s
}

func (s *ShipSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("ship")
	for id, c := range comps {
		ship := c.(*components.ShipComponent)

		c, _ = s.getter.GetComponent(id, "movement")
		movement := c.(*components.MovementComponent)

		s.moveShip(ship, movement)
	}

	return nil
}

func (s *ShipSystem) moveShip(ship *components.ShipComponent, movement *components.MovementComponent) {

	if ship.AvailableID != "" && s.getter.HasComponents(ship.AvailableID, "control") {
		c, _ := s.getter.GetComponent(ship.AvailableID, "control")
		control := c.(*components.ControlComponent)

		movement.Direction = getMoveVector(control)
	}

}

func getMoveVector(control *components.ControlComponent) geometry.Vec3 {
	vec := geometry.GetZeroVector()
	if control.MoveRight {
		vec.Add(geometry.Vec3{X: 1, Y: 1, Z: 0})
	}
	if control.MoveLeft {
		vec.Add(geometry.Vec3{X: 1, Y: -1, Z: 0})
	}
	return vec.Normalize()
}
