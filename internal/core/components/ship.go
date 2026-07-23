package components

type ShipComponent struct {
	AvailableID   string
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
