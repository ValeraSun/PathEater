package components

type NavigationEntityComponent struct {
}

func (*NavigationEntityComponent) Type() string {
	return "navigationEntity"
}

func NewNavigationEntityComponent() *NavigationEntityComponent {
	return &NavigationEntityComponent{}
}
