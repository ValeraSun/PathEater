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
	damage                int
}

func (*AttackComponent) Type() string {
	return "attack"
}

func NewAttackComponent(damage, cooldown, distant float64, hitbox geometry.Collider) *AttackComponent {
	return &AttackComponent{
		damage:         int(damage),
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

func (c *AttackComponent) TryAttack(id types.Entity, pos, distant geometry.Vec3) (*events.AttackEvent, bool) {
	isAttack := distant.Length() < 2 && c.currentAttackCooldown == 0
	if isAttack {
		c.hitbox.ChangeCenter(pos.Add(distant.Normalize().Scale(c.attackRange)))
		c.currentAttackCooldown = c.attackCooldown

		e := events.NewAttackEvent(
			c.damage,
			c.hitbox,
			id,
		)
		return e, isAttack
	}
	return nil, isAttack
}
