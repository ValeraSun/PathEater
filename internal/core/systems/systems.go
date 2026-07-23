package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type entityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
	AddEntityByID(entity types.Entity, components ...types.Component) error
}

type closer interface {
	Close()
}

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}

type entitiesRemover interface {
	RemoveEntity(entity types.Entity)
}

type Broadcaster interface {
	SendEntityCreate(ecs.EntityInfo) error
	SendEntityUpdate(ecs.EntityInfo) error
	SendEntityDelete(ecs.EntityInfo) error
	SendGameOverState(ecs.GameOverInfo) error
}

type subscriber interface {
	Subscribe(eventType string, handler events.EventHandler) (func(), error)
}

type publisher interface {
	Publish(events.Event) error
}

type entityRemover interface {
	RemoveEntity(entity types.Entity)
}
