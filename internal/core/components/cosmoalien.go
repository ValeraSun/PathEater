package components

type CosmoAlienComponent struct{}

func (*CosmoAlienComponent) Type() string {
	return "cosmoAlien"
}

func NewCosmoAlienComponent() *CosmoAlienComponent {
	return &CosmoAlienComponent{}
}
