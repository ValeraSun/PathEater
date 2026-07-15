package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type ExternalVelocityComponent struct {
	Direction geometry.Vec3
	Speed     float64
}

func (*ExternalVelocityComponent) Type() string {
	return "externalVelocity"
}

func NewExternalVelocityComponent() *ExternalVelocityComponent {
	return &ExternalVelocityComponent{}
}
