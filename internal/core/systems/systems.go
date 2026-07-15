package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type entityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
}

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}
type Broadcaster interface {
	SendEntityCreate(ecs.EntityCreateInfo) error
	SendEntityUpdate(ecs.EntityUpdateInfo) error
	SendEntityDelete(ecs.EntityUpdateInfo) error
	CreateCameraForPlayer(id string, idEntity types.Entity) error
}

type subscriber interface {
	Subscribe(eventType string, handler events.EventHandler) (func(), error)
}

type publisher interface {
	Subscribe(eventType string, handler events.EventHandler) (func(), error)
}
