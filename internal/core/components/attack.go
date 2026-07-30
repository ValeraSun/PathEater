package components

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type AttackComponent struct {
	currentAttackCooldown float32
	attackCooldown        float32
	attackRange           float64
	hitbox                geometry.Collider
	damage                float32
}

func (*AttackComponent) Type() string {
	return "attack"
}

func NewAttackComponent(damage, cooldown, distant float64, hitbox geometry.Collider) *AttackComponent {
	return &AttackComponent{
		damage:         float32(damage),
		attackCooldown: float32(cooldown),
		attackRange:    distant,
		hitbox:         hitbox,
	}
}

func (c *AttackComponent) ReduceCooldown(dt float32) {
	if c.currentAttackCooldown-dt > 0 {
		c.currentAttackCooldown -= dt
	} else {
		c.currentAttackCooldown = 0
	}
}

func (c *AttackComponent) Attack(id types.Entity, pos, direction geometry.Vec3) (*events.AttackEvent, bool) {
	if c.currentAttackCooldown == 0 {
		c.hitbox.ChangeCenter(pos.Add(direction.Scale(c.attackRange)))
		c.currentAttackCooldown = c.attackCooldown

		e := events.NewAttackEvent(
			c.damage,
			c.hitbox,
			id,
		)

		return e, true
	}
	return nil, false
}
