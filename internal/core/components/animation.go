package components

import "github.com/ValeraSun/PathEater/internal/core/events"

type AnimationComponent struct {
	state    string
	cooldown float32
}

func (*AnimationComponent) Type() string {
	return "animation"
}

func NewAnimationComponent(state string) *AnimationComponent {
	return &AnimationComponent{
		state,
		0,
	}
}

func (c *AnimationComponent) ResetCooldown(dt float32) {
	if c.cooldown-dt > 0 {
		c.cooldown -= dt
	} else {
		c.cooldown = 0
		c.state = "base"
	}
}

func (c *AnimationComponent) ChangeState(event *events.AnimationEvent) {
	c.state = event.State
	c.cooldown = event.Cooldown
}

func (c *AnimationComponent) CurrentAnimation() string {
	return c.state
}
