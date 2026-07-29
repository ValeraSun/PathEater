package systems

import (
	"time"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/sendler"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	eventQueueSize = 100
)

type entityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
	AddEntityByID(entity types.Entity, components ...types.Component) error
}

type closer interface {
	Close()
}

type timer interface {
	GetRemainingTime() time.Duration
}

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
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

type entitySendler interface {
	Send(types.Entity, func(sendler.EntityInfo) error) error
}
type broadcasterFunc interface {
	SendEntityCreate(entityInfo sendler.EntityInfo) error
	SendEntityUpdate(entityInfo sendler.EntityInfo) error
	SendEntityDelete(entityInfo sendler.EntityInfo) error
	SendGameOverState(sendler.GameOverInfo) error
	SendTime(sendler.TimeInfo) error
}

type broadcaster interface {
	entitySendler
	broadcasterFunc
}

func getShip(getter componentsGetter) (types.Entity, *components.ShipComponent) {
	ships := getter.GetEntitiesByComponent("ship")

	for id, s := range ships {
		return id, s.(*components.ShipComponent)
	}
	return "", nil
}
