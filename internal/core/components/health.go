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

func (hp *HealthComponent) Damage(damage int) bool {
	hp.Health = hp.Health - damage
	if hp.Health <= 0 {
		return true
	}
	return false
}
