package components

type CosmoAlienComponent struct {
	Visible bool
	Died    bool
}

func (*CosmoAlienComponent) Type() string {
	return "cosmoAlien"
}

func NewCosmoAlienComponent() *CosmoAlienComponent {
	return &CosmoAlienComponent{
		Visible: false,
		Died:    false,
	}
}
