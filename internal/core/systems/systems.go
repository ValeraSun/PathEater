package systems

import "github.com/ValeraSun/PathEater/internal/core/types"

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
}
