package entities

import (
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type entityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
	AddEntityByID(entity types.Entity, components ...types.Component) error
}

// func NewBox(width, depth, height float64, adder entityAdder) types.Entity {
// 	e, _ := adder.AddEntity(
// 		components.NewTransformComponent(geometry.GetZeroVector(), geometry.GetZeroVector()),
// 		components.NewColliderComponent(geometry.NewBoxCollider(
// 			geometry.GetZeroVector(),
// 			geometry.Vec3{X: width / 2, Y: height / 2, Z: depth / 2})),
// 	)

// 	return e
// }

func NewPlayer(adder entityAdder, clientID string) types.Entity {
	adder.AddEntityByID(
		types.Entity(clientID),
		components.NewUpdateComponent(),
		components.NewControlComponent(clientID),
		components.NewMovementComponent(10),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewHealthComponent(100),
		components.NewColliderComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{},
			geometry.Vec3{Y: 1},
			1.0, 1.0)),
		components.NewMovableComponent(),
		components.NewTransformComponent(
			geometry.Vec3{
				X: 2,
				Y: 1,
				Z: 0,
			},
			geometry.GetZeroVector(),
		),
	)
	return types.Entity(clientID)

}

func NewShip(adder entityAdder) types.Entity {
	e, _ := adder.AddEntity(
		components.NewShipComponent(),
		components.NewHealthComponent(100),
	)
	return e
}

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}

type Broadcaster interface {
	SendEntityCreate(ecs.EntityInfo) error
	SendEntityUpdate(ecs.EntityInfo) error
	SendEntityDelete(ecs.EntityInfo) error
}

type playerData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
	Health   int
}

type shipData struct {
	baggageStatus int
	health        int
}

type Sendler func(ecs.EntityInfo) error

func IsShip(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "ship") && getter.HasComponents(id, "health")
}
func SendShip(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	c, ok := getter.GetComponent(id, "ship")
	ship := c.(*components.ShipComponent)

	if !ok {
		return errors.New("Не найден необхадимый компнонент")
	}

	c, ok = getter.GetComponent(id, "health")
	hp := c.(*components.HealthComponent)

	if !ok {
		return errors.New("Не найден необхадимый компнонент")
	}

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "ship",
		Data: shipData{
			baggageStatus: ship.StatusBag,
			health:        hp.Health,
		}})

	return nil
}

func IsPlayer(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "transform") && getter.HasComponents(id, "control") && getter.HasComponents("health")
}

func SendPlayer(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	c, _ := getter.GetComponent(id, "transform")
	transform := c.(*components.TransformComponent)

	c, _ = getter.GetComponent(id, "health")
	hp := c.(*components.HealthComponent)

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "player",
		Data: playerData{
			Position: transform.Position,
			Rotation: transform.Direction,
			Health:   hp.Health,
		},
	},
	)

	return nil
}
