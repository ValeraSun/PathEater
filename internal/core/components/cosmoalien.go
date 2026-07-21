package components

type CosmoAlienComponent struct{
	died bool
}

func (*CosmoAlienComponent) Type() string {
	return "cosmoAlien"
}

func NewCosmoAlienComponent() *CosmoAlienComponent {
	return &CosmoAlienComponent{
		died: false,
	}
}
