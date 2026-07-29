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
