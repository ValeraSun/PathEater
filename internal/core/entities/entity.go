package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type EntityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
}

func NewPlayer(adder EntityAdder, id string) types.Entity {
	transform := components.NewTransformComponent()
	velocity := components.NewVelocityComponent()
	control := components.NewControlComponent(id)
	entity, _ := adder.AddEntity(transform, velocity, control) //добавить обработку

	return entity
}
