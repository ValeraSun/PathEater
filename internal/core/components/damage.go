package components

type DamageAreaComponent struct {
	Health    int
	MaxHealth int
}

func NewDamageAreaComponent(maxHealth int) *DamageAreaComponent {
	return &DamageAreaComponent{
		Health:    maxHealth,
		MaxHealth: maxHealth,
	}
}

func (*DamageAreaComponent) Type() string { return "damageArea" }
