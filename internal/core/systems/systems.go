package systems

import "github.com/ValeraSun/PathEater/internal/core/types"

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}

type entityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
}
