package components

type ShipComponent struct {
	StatusBag int
}

func NewShipComponent() *ShipComponent {
	return &ShipComponent{
		StatusBag: 100,
	}
}
func (*ShipComponent) Type() string { return "ship" }
