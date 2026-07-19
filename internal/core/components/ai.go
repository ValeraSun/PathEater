package components

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type AIComponent struct {
	MoveFront   bool
	Direction   geometry.Vec3
	WantPostion geometry.Vec3
}

func (*AIComponent) Type() string {
	return "ai"
}

func NewAIComponent() *AIComponent {
	return &AIComponent{}
}

func (c *AIComponent) GetInputVector() geometry.Vec3 {
	input := geometry.Vec3{}

	return input.Normalize()
}
