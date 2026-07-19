package components

type ShipComponent struct {
	BaggageStatus int
}

func (*ShipComponent) Type() string {
	return "ship"
}

func NewShipComponent(baggageStatus int) *ShipComponent {
	return &ShipComponent{
		BaggageStatus: baggageStatus,
	}
}
