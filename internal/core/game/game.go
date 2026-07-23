package game

import (
	"github.com/ValeraSun/PathEater/internal/config"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/systems"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type subscriber interface {
	Subscribe(eventType string, handler events.EventHandler) (func(), error)
}

type closer interface {
	Close()
}

type systemAdder interface {
	AddSystem(types.System)
	AddEntity(components ...types.Component) (types.Entity, error)
	AddEntityByID(entity types.Entity, components ...types.Component) error
}

type entityRemover interface {
	RemoveEntity(entity types.Entity)
}

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}

func CreateGame(broadcaster ecs.Broadcaster) *ecs.World {
	w := ecs.CreateWorld(broadcaster)
	initSystems(w, w, w, w.EventBus, w.Broadcaster, w.EventBus, w)
	config.CreateWalls(w)

	createEntities(w)

	go w.EventBus.ProcessEvents()
	go ecs.HandleWorld(w)
	return w
}

func initSystems(adder systemAdder, remover entityRemover, getter componentsGetter, publisher transfer.EventPublisher, broadcaster ecs.Broadcaster, subscriber subscriber, closer closer) {
	adder.AddSystem(systems.NewVisionSystem(getter))
	adder.AddSystem(systems.NewAISystem(getter))
	adder.AddSystem(systems.NewControlSystem(getter, subscriber))
	adder.AddSystem(systems.NewMovementSystem(getter))
	adder.AddSystem(systems.NewVelocitySystem(getter))
	adder.AddSystem(systems.NewTransformSystem(getter))
	adder.AddSystem(systems.NewCollisionSystem(getter, publisher))
	adder.AddSystem(systems.NewCollisionsSystem(getter, publisher, subscriber))
	adder.AddSystem(systems.NewRenderSystem(getter, broadcaster))
	adder.AddSystem(systems.NewCreateSystem(adder, getter, broadcaster, subscriber))
	adder.AddSystem(systems.NewDeleteSystem(getter, broadcaster, remover, subscriber))
	adder.AddSystem(systems.NewWeaponSystem(getter, publisher, subscriber))
	adder.AddSystem(systems.NewShootSystem(getter))
	adder.AddSystem(systems.NewAsteroidSystem(getter, publisher, subscriber))
	adder.AddSystem(systems.NewShipSystem(getter, subscriber))
	adder.AddSystem(systems.NewHealthSystem(getter, publisher, subscriber))
	adder.AddSystem(systems.NewCosmoAlienSystem(getter, publisher))
	adder.AddSystem(systems.NewGameOverSystem(getter, subscriber, broadcaster, closer))
	adder.AddSystem(systems.NewDeathSystem(getter, publisher, subscriber))
}

func createEntities(world *ecs.World) {
	world.EventBus.Publish(events.NewCreateShipEvent())
	world.EventBus.Publish(events.NewCreateAlienEvent(geometry.Vec3{X: 3, Y: 2, Z: -1}))
}
