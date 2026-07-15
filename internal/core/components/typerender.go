package components

type TypeRenderComponent struct {
	TypeRender string
}

func NewTypeRenderComponent(typeRender string) *TypeRenderComponent {
	return &TypeRenderComponent{
		TypeRender: typeRender,
	}
}

func (*TypeRenderComponent) Type() string { return "typeRender" }
