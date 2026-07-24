package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
)

const (
	forwardSpeed = 30.0
	verticalSpeed = 20.0
)

type ShipSystem struct {
	getter componentsGetter
}

func NewShipSystem(getter componentsGetter) *ShipSystem {
	return &ShipSystem{
		getter: getter,
	}
}

func (s *ShipSystem) Update(dt float32) error {
	s.handleShipControl()
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

	if !s.getter.HasComponents(shipID, "externalVelocity") {
		return
	}
	extComp, _ := s.getter.GetComponent(shipID, "externalVelocity")
	ext := extComp.(*components.ExternalVelocityComponent)

	ext.Direction.X = forwardSpeed

	if control.MoveFront {
		ext.Direction.Y = verticalSpeed
	} else if control.MoveBack {
		ext.Direction.Y = -verticalSpeed
	} else {
		ext.Direction.Y = 0
	}
	ext.Direction.Z = 0
}

func (s *ShipSystem) getShipID() string {
	shipEntities := s.getter.GetEntitiesByComponent("ship")
	for id := range shipEntities {
		return id
	}
	return ""
}