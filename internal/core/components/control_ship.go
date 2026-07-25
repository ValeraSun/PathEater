package components

type ControlShipComponent struct {
	IsControling bool
}

func NewControlShipComponent() *ControlShipComponent {
	return &ControlShipComponent{
		IsControling: false,
	}
}

func (*ControlShipComponent) Type() string { return "controlShip" }
