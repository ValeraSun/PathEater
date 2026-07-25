package game

import (
	"time"

	"github.com/ValeraSun/PathEater/internal/config"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/systems"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type Subscriber interface {
	Subscribe(eventType string, handler events.EventHandler) (func(), error)
}

type closer interface {
	Close()
}

type timer interface {
	GetRemainingTime() time.Duration
}

type systemAdder interface {
	AddSystem(types.System)
	AddEntity(components ...types.Component) (types.Entity, error)
	AddEntityByID(entity types.Entity, components ...types.Component) error
}

type entityRemover interface {
	RemoveEntity(entity types.Entity)
}

type componentsworld interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}

func CreateGame(broadcaster ecs.Broadcaster) *ecs.World {
	w := ecs.CreateWorld(broadcaster)
	initSystems(w)
	config.CreateWalls(w)

	createEntities(w)

	go w.EventBus.ProcessEvents()
	go ecs.HandleWorld(w)
	w.Timer.StartTimer(99999*time.Second, func() { timerIsOver(w.EventBus) })
	return w
}

func initSystems(world *ecs.World) {
	world.AddSystem(systems.NewVisionSystem(world))
	world.AddSystem(systems.NewAISystem(world, world.EventBus))
	world.AddSystem(systems.NewAttackSystem(world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewControlSystem(world, world.EventBus))
	world.AddSystem(systems.NewNavigationDisplaySystem(world))
	world.AddSystem(systems.NewMovementSystem(world))
	world.AddSystem(systems.NewExternalVelocitySystem(world))
	world.AddSystem(systems.NewVelocitySystem(world))
	world.AddSystem(systems.NewTransformSystem(world))
	world.AddSystem(systems.NewCollisionSystem(world, world.EventBus))
	world.AddSystem(systems.NewCollisionsSystem(world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewRenderSystem(world, world.Broadcaster))
	world.AddSystem(systems.NewCreateSystem(world, world, world.Broadcaster, world.EventBus))
	world.AddSystem(systems.NewDeleteSystem(world, world.Broadcaster, world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewWeaponSystem(world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewShootSystem(world))
	world.AddSystem(systems.NewAsteroidSystem(world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewHealthSystem(world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewCosmoAlienSystem(world, world.EventBus))
	world.AddSystem(systems.NewGameOverSystem(world, world.EventBus, world.Broadcaster, world))
	world.AddSystem(systems.NewDeadSystem(world, world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewDeathSystem(world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewTimerSystem(world, world.Broadcaster, &world.Timer))
}

func createEntities(world *ecs.World) {
	world.EventBus.Publish(events.NewCreateShipEvent())
	//world.EventBus.Publish(events.NewCreateAlienEvent(geometry.Vec3{X: 3, Y: 2, Z: -1}))
	world.EventBus.Publish(events.NewMeteoriteZoneEvent())
}

func timerIsOver(publisher transfer.EventPublisher) {
	publisher.Publish(events.NewGameOverEvent(true))
}
