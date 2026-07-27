package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func NewTerminal(adder entityAdder) types.Entity {
	e, _ := adder.AddEntity(
		components.NewInteractableComponent(
			geometry.Vec3{X: 39.85, Y: 1.075, Z: -9},
			geometry.Vec3{X: 0.5, Y: 10, Z: 3},
			components.InteractTerminal,
		),
	)
	return e
}
