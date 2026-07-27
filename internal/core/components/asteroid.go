package components

type AsteroidComponent struct {
	Destroyed bool
}

func (*AsteroidComponent) Type() string {
	return "asteroid"
}

func NewAsteroidComponent() *AsteroidComponent {
	return &AsteroidComponent{
		Destroyed: false,
	}
}
