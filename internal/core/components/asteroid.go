package components

type AsteroidComponent struct{}

func (*AsteroidComponent) Type() string {
	return "asteroid"
}

func NewAsteroidComponent() *AsteroidComponent {
	return &AsteroidComponent{}
}
