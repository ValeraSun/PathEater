package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

const (
	forwardSpeed = 30.0
	verticalSpeed = 20.0
)

type ShipSystem struct {
	getter     componentsGetter
}

func NewShipSystem(getter componentsGetter, subscriber subscriber) *ShipSystem {
	return &ShipSystem{
		getter: getter,
	}
}

func (s *ShipSystem) Update(dt float32) error {
	s.handleShipControl()

	entities := s.getter.GetEntitiesByComponent("navigationEntity")
	for id, _ := range entities {
		c, _ := s.getter.GetComponent(id, "movement")
		movement := c.(*components.MovementComponent)
		movement.Position.X += movement.Velocity.X * float64(dt)
		movement.Position.Y += movement.Velocity.Y * float64(dt)
		movement.Position.Z += movement.Velocity.Z * float64(dt)
	}

	return nil
}

func (s *ShipSystem) handleShipControl() {
	shipID := s.getShipID()
	if shipID == "" {
		return
	}

	if !s.getter.HasComponents(shipID, "control") {
		return
	}
	ctrlComp, _ := s.getter.GetComponent(shipID, "control")
	control := ctrlComp.(*components.ControlComponent)

	movComp, ok := s.getter.GetComponent(shipID, "movement")
	if !ok {
		return
	}
	movement := movComp.(*components.MovementComponent)

	movement.Velocity.X = forwardSpeed

	if control.MoveFront {
		movement.Velocity.Y = verticalSpeed
	} else if control.MoveBack {
		movement.Velocity.Y = -verticalSpeed
	} else {
		movement.Velocity.Y = 0
	}
}

func (s *ShipSystem) getShipID() string {
	shipEntities := s.getter.GetEntitiesByComponent("ship")
	for id := range shipEntities {
		return id
	}
	return ""
}