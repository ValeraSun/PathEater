package components

type CosmoAlienComponent struct {
	Died bool
}

func (*CosmoAlienComponent) Type() string {
	return "cosmoAlien"
}

func NewCosmoAlienComponent() *CosmoAlienComponent {
	return &CosmoAlienComponent{
		Died: false,
	}
}
