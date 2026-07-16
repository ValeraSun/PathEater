package game

import (
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/systems"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type subscriber interface {
	Subscribe(eventType string, handler events.EventHandler) (func(), error)
}
type systemAdder interface {
	AddSystem(types.System)
	AddEntity(components ...types.Component) (types.Entity, error)
}

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}

func CreateGame(broadcaster ecs.Broadcaster) *ecs.World {
	w := ecs.CreateWorld(broadcaster)
	//eb := *w.EventBus
	initSystems(w, w, w.Broadcaster, w.EventBus)

	go w.EventBus.ProcessEvents()
	go ecs.HandleWorld(w)
	return w
}

func initSystems(adder systemAdder, getter componentsGetter, broadcaster ecs.Broadcaster, subscriber subscriber) {
	adder.AddSystem(systems.NewControlSystem(getter, subscriber))
	adder.AddSystem(systems.NewMovementSystem(getter))
	adder.AddSystem(systems.NewVelocitySystem(getter))
	adder.AddSystem(systems.NewTransformSystem(getter))
	adder.AddSystem(systems.NewCollisionSystem(getter))
	adder.AddSystem(systems.NewRenderSystem(getter, broadcaster))
	adder.AddSystem(systems.NewCreateSystem(adder, getter, broadcaster, subscriber))

}
