package components

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type AIComponent struct {
	MoveFront      bool
	Direction      geometry.Vec3
	AttackCooldown float32
}

func (*AIComponent) Type() string {
	return "ai"
}

func NewAIComponent() *AIComponent {
	return &AIComponent{}
}
