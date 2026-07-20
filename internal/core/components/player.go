package components

type PlayerComponent struct{}

func (*PlayerComponent) Type() string {
	return "player"
}

func NewPlayerComponent() *PlayerComponent {
	return &PlayerComponent{}
}
