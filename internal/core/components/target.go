package components

type TargetComponent struct{}

func (*TargetComponent) Type() string {
	return "target"
}

func NewTargetComponent() *TargetComponent {
	return &TargetComponent{}
}
