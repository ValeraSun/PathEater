package components

type ShipComponent struct {
	BaggageStatus int
	Speed         float64
}

func (*ShipComponent) Type() string {
	return "ship"
}

func NewShipComponent(baggageStatus int, speed float64) *ShipComponent {
	return &ShipComponent{
		BaggageStatus: baggageStatus,
		Speed:         speed,
	}
}
