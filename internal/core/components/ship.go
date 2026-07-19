package components

type ShipComponent struct {
	BaggageStatus int
}

func (*ShipComponent) Type() string {
	return "ship"
}

func NewShipComponent() *ShipComponent {
	return &ShipComponent{
		BaggageStatus: 10,
	}
}
