package components

type AnimationComponent struct {
	state string
}

func (*AnimationComponent) Type() string {
	return "animation"
}

func NewAnimationComponent() *AnimationComponent {
	return &AnimationComponent{}
}
