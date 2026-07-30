package components

type RayAttackComponent struct {
	currentAttackCooldown float32
	attackCooldown        float32
	Damage                float32
}

func (*RayAttackComponent) Type() string { return "rayAttack" }

func NewRayAttackComponent(attackCooldown float32, damage float32) *RayAttackComponent {
	return &RayAttackComponent{
		attackCooldown: attackCooldown,
		Damage:         damage,
	}
}
func (c *RayAttackComponent) ReduseCooldown(dt float32) {
	if c.currentAttackCooldown-dt > 0 {
		c.currentAttackCooldown -= dt
	} else {
		c.currentAttackCooldown = 0
	}
}

func (c *RayAttackComponent) ResetCooldown() {
	c.currentAttackCooldown = c.attackCooldown
}

func (c *RayAttackComponent) CanAttack() bool {
	return c.currentAttackCooldown == 0
}
