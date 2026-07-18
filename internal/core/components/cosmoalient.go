package components

type CosmoAlientComponent struct{}

func (*CosmoAlientComponent) Type() string {
	return "cosmoAlient"
}

func NewCosmoAlientComponent() *CosmoAlientComponent {
	return &CosmoAlientComponent{}
}
