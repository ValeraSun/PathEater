package components

type BaggageComponent struct{}

func (*BaggageComponent) Type() string {
	return "baggage"
}

func NewBaggageComponent() *BaggageComponent {
	return &BaggageComponent{}
}
