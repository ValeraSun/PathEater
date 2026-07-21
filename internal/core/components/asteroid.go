package components

type AsteroidComponent struct {
	Visible   bool
	OnField   bool
	Destroyed bool
}

func (*AsteroidComponent) Type() string {
	return "asteroid"
}

func NewAsteroidComponent() *AsteroidComponent {
	return &AsteroidComponent{
		Visible:   false,
		OnField:   true,
		Destroyed: false,
	}
}
