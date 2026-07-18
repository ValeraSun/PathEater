package components

type UpdateComponent struct {
}

func NewUpdateComponent() *UpdateComponent {
	return &UpdateComponent{}
}

func (*UpdateComponent) Type() string { return "update" }
