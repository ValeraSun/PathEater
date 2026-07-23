package components

type PlayerComponent struct {
	Died bool
}

func (*PlayerComponent) Type() string {
	return "player"
}

func NewPlayerComponent() *PlayerComponent {
	return &PlayerComponent{
		Died: false,
	}
}
