package components

type HealthComponent struct {
	Health    int
	MaxHealth int
}

func NewHealthComponent(maxHealth int) *HealthComponent {
	return &HealthComponent{
		Health:    maxHealth,
		MaxHealth: maxHealth,
	}
}

func (*HealthComponent) Type() string { return "health" }
