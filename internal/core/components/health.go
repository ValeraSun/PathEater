package components

type HealthComponent struct {
	Health    float32
	MaxHealth float32
}

func NewHealthComponent(maxHealth float32) *HealthComponent {
	return &HealthComponent{
		Health:    maxHealth,
		MaxHealth: maxHealth,
	}
}

func (*HealthComponent) Type() string { return "health" }

func (hp *HealthComponent) Damage(damage float32) bool {
	hp.Health = hp.Health - damage
	if hp.Health <= 0 {
		return true
	}
	return false
}
