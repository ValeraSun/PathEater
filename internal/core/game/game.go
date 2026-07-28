package game

import (
	"time"

	"github.com/ValeraSun/PathEater/internal/config"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/sendler"
	"github.com/ValeraSun/PathEater/internal/core/systems"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	timeGame = 9999
)

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

func CreateGame(broadcasterFunc broadcasterFunc) *ecs.World {
	w := ecs.CreateWorld(broadcasterFunc)
	initSystems(w)
	config.CreateWalls(w)
	config.CreateDoors(w)
	config.CreateRooms(w)
	createEntities(w)

	go w.EventBus.ProcessEvents()
	go ecs.HandleWorld(w)
	w.Timer.StartTimer(timeGame*time.Second, func() { timerIsOver(w.EventBus) })
	return w
}

func initSystems(world *ecs.World) {
	world.AddSystem(systems.NewVisionSystem(world))
	world.AddSystem(systems.NewAISystem(world, world.EventBus))
	world.AddSystem(systems.NewAttackSystem(world, world.EventBus, world.EventBus)) //alien ai

	world.AddSystem(systems.NewControlSystem(world, world.EventBus))
	world.AddSystem(systems.NewInteractionSystem(world, world.EventBus))
	world.AddSystem(systems.NewRayAttackComponent(world, world.EventBus)) //игрок

	world.AddSystem(systems.NewBoardingSystem(world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewAsteroidSystem(world, world.EventBus, world.EventBus))
	world.AddSystem(systems.NewCosmoAlienSystem(world, world.EventBus))
	//	world.AddSystem(systems.NewWeaponSystem(world, world.EventBus, world.EventBus))
	//	world.AddSystem(systems.NewShootSystem(world))
	world.AddSystem(systems.NewNavigationDisplaySystem(world, world.EventBus))
	world.AddSystem(systems.NewCollisionsSystem(world, world.EventBus, world.EventBus)) //терминал

	world.AddSystem(systems.NewMovementSystem(world))
	world.AddSystem(systems.NewExternalVelocitySystem(world))
	world.AddSystem(systems.NewVelocitySystem(world))
	world.AddSystem(systems.NewTransformSystem(world))
	world.AddSystem(systems.NewColliderMoveSystem(world))
	world.AddSystem(systems.NewCollisionSystem(world, world.EventBus)) //передвижение

	world.AddSystem(systems.NewBreakdownSystem(world, world.EventBus, world.EventBus, config.GetExternalWalls()))
	world.AddSystem(systems.NewTimerSystem(world, world.Broadcaster, &world.Timer))
	world.AddSystem(systems.NewDoorsSystem(world, world.EventBus))
	world.AddSystem(systems.NewRoomsSystem(world))
	world.AddSystem(systems.NewVacuumSystem(world, world.EventBus))
	world.AddSystem(systems.NewOxygenSystem(world, world.EventBus))
	world.AddSystem(systems.NewWeaponSystem(world, world.EventBus, world.EventBus))

	world.AddSystem(systems.NewHealthSystem(world, world.EventBus, world.EventBus)) // здоровье

	world.AddSystem(systems.NewTimerSystem(world, world.Broadcaster, &world.Timer))
	world.AddSystem(systems.NewGameOverSystem(world, world.EventBus, world.Broadcaster, world)) // конец игры

	world.AddSystem(systems.NewRenderSystem(world, world.Broadcaster))
	world.AddSystem(systems.NewCreateSystem(world, world, world.EventBus, world.Broadcaster))
	world.AddSystem(systems.NewDeleteSystem(world, world, world.EventBus, world.EventBus, world.Broadcaster)) // база

}

func createEntities(world *ecs.World) {
	entities.NewTerminal(world)
	world.EventBus.Publish(events.NewCreateShipEvent())
	entities.InitWeapon(world)
	//world.EventBus.Publish(events.NewCreateAlienEvent(geometry.Vec3{X: 3, Y: 2, Z: -1}))
	//world.EventBus.Publish(events.NewCreateAlienEvent(geometry.Vec3{X: 3, Y: 2, Z: -1}))
	world.EventBus.Publish(events.NewMeteoriteZoneEvent())
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
	world.EventBus.Publish(events.NewBreakdownEvent(true))
}

func timerIsOver(publisher transfer.EventPublisher) {
	publisher.Publish(events.NewGameOverEvent(true))
}
