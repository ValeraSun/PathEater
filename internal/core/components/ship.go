package components

type ShipComponent struct{}

func (*ShipComponent) Type() string {
	return "ship"
}

func NewShipComponent() *ShipComponent {
	return &ShipComponent{}
}
