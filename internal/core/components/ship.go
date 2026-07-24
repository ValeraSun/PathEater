package components

import "github.com/ValeraSun/PathEater/internal/core/types"

type ShipComponent struct {
	AvailableID   types.Entity
	BaggageStatus int
	Speed         float64
}

func (*ShipComponent) Type() string {
	return "ship"
}

func NewShipComponent(baggageStatus int, speed float64) *ShipComponent {
	return &ShipComponent{
		AvailableID:   "",
		BaggageStatus: baggageStatus,
		Speed:         speed,
	}
}
