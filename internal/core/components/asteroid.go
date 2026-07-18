package components

type AsteroidComponent struct {
	Visible   bool
	Destroyed bool
}

func (*AsteroidComponent) Type() string {
	return "asteroid"
}

func NewAsteroidComponent() *AsteroidComponent {
	return &AsteroidComponent{
		Visible:   false,
		Destroyed: false,
	}
}
