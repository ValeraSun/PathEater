package components

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type ShipComponent struct {
	AvailableID   types.Entity
	BaggageStatus int
	Speed         float64
}

func (*ShipComponent) Type() string {
	return "ship"
}

func NewShipComponent(baggageStatus int) *ShipComponent {
	return &ShipComponent{
		AvailableID:   "",
		BaggageStatus: baggageStatus,
	}
}

type componentsGetter interface {
	GetEntitiesByComponent(string) map[types.Entity]types.Component
}

func GetShip(getter componentsGetter) types.Entity {
	ships := getter.GetEntitiesByComponent("ship")

	var ship types.Entity

	for id := range ships {
		ship = id
	}

	return ship
}
