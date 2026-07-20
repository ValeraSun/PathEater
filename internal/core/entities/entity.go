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
		//components.NewUpdateComponent(),
		components.NewControlComponent(clientID),
		components.NewMovementComponent(10, geometry.GetZeroVector()),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewHealthComponent(100),
		components.NewColliderComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{X: 3, Y: 1},
			geometry.Vec3{Y: 1},
			1, 0.7)),
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
		components.NewTransformComponent(
			geometry.GetZeroVector(),
			geometry.GetZeroVector(),
		),
		components.NewHealthComponent(100),
		components.NewControlComponent(""),
		components.NewShipComponent(10, 5),
		components.NewWeaponComponent(geometry.GetZeroVector(), 10, 30),
		components.NewColliderComponent(geometry.NewTriangleCollider(
			geometry.Vec3{X: 0, Y: 20, Z: 0},
			geometry.Vec3{X: 20, Y: -20, Z: 0},
			geometry.Vec3{X: -20, Y: -20, Z: 0},
		)),
	)
	return e
}

func NewAsteroid(adder entityAdder, pos geometry.Vec3, radius float64, dir geometry.Vec3, speed float64) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			pos,
			dir,
		),
		components.NewMovementComponent(speed, dir),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewAsteroidComponent(),
		components.NewColliderComponent(geometry.NewCircleCollider(
			pos,
			radius,
		)),
	)
	return e
}

func NewCosmoAlien(adder entityAdder, pos geometry.Vec3, shipID types.Entity) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			pos,
			geometry.GetZeroVector(),
		),
		components.NewMovementComponent(10, geometry.GetZeroVector()),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewStalkerComponent(shipID),
		components.NewCosmoAlienComponent(),
		components.NewColliderComponent(geometry.NewCircleCollider(
			pos,
			30,
		)),
		components.NewMovableComponent(),
	)
	return e
}

func NewBullet(adder entityAdder, pos, dir geometry.Vec3) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			pos,
			dir,
		),
		components.NewBulletComponent(),
		components.NewMovementComponent(10, dir),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewColliderComponent(geometry.NewCircleCollider(
			pos,
			2,
		)),
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
	Collider        geometry.Collider `json:"collider"`
	WeaponDirection geometry.Vec3     `json:"weapon_direction"`
	ShootSuccess    bool              `json:"shoot_success"`
	BaggageStatus   int               `json:"baggage_status"`
	Health          int               `json:"health"`
}

type asteroidData struct {
	Position  geometry.Vec3 `json:"position"`
	Radius    float64       `json:"radius"`
	Destroyed bool          `json:"destroyed"`
}

type cosmoAlienData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
	Health   int           `json:"health"`
	Died     bool          `json:"died"`
}

type bulletData struct {
	Position geometry.Vec3 `json:"position"`
	Success  bool          `json:"success"`
}

type Sendler func(ecs.EntityInfo) error

func IsPlayer(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "transform") && getter.HasComponents(id, "control") && getter.HasComponents(id, "health")
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

func IsShip(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "ship")
}

func SendShip(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	c, ok := getter.GetComponent(id, "ship")
	ship := c.(*components.ShipComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	c, ok = getter.GetComponent(id, "health")
	hp := c.(*components.HealthComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	c, ok = getter.GetComponent(id, "collider")
	col := c.(*components.ColliderComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	c, ok = getter.GetComponent(id, "weapon")
	weap := c.(*components.WeaponComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "ship",
		Data: shipData{
			Collider:        col.Collider,
			WeaponDirection: weap.Direction,
			ShootSuccess:    weap.ShootSuccess,
			BaggageStatus:   ship.BaggageStatus,
			Health:          hp.Health,
		}})

	return nil
}

func IsAsteroid(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "asteroid")
}

func IsVisible(id types.Entity, getter componentsGetter) bool {
	c, _ := getter.GetComponent(id, "asteroid")
	aster := c.(*components.AsteroidComponent)
	return aster.Visible
}

func SendAsteroid(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	c, ok := getter.GetComponent(id, "asteroid")
	aster := c.(*components.AsteroidComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	c, ok = getter.GetComponent(id, "transform")
	transform := c.(*components.TransformComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	c, ok = getter.GetComponent(id, "collider")
	col := c.(*components.ColliderComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "asteroid",
		Data: asteroidData{
			Position:  transform.Position,
			Radius:    col.Collider.(*geometry.CircleCollider).Radius,
			Destroyed: aster.Destroyed,
		}})

	return nil
}

func IsCosmoAlien(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "cosmoAlien")
}

func SendCosmoAlien(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	c, ok := getter.GetComponent(id, "cosmoAlien")
	alien := c.(*components.CosmoAlienComponent)

	if !ok {
		return errors.New("Не найден необходимый компонент")
	}

	c, ok = getter.GetComponent(id, "transform")
	transform := c.(*components.TransformComponent)

	if !ok {
		return errors.New("Не найден необходимый компонент")
	}

	c, ok = getter.GetComponent(id, "health")
	hp := c.(*components.HealthComponent)

	if !ok {
		return errors.New("Не найден необходимый компонент")
	}

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "alien",
		Data: cosmoAlienData{
			Position: transform.Position,
			Rotation: transform.Direction,
			Health:   hp.Health,
			Died:     alien.Died
		}})

	return nil
}

func IsBullet(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "bullet")
}

func SendBullet(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	c, ok := getter.GetComponent(id, "transform")
	transform := c.(*components.TransformComponent)

	if !ok {
		return errors.New("Не найден необходимый компонент")
	}

	if !ok {
		return errors.New("Не найден необходимый компонент")
	}

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "bullet",
		Data: bulletData{
			Position: transform.Position,
		}})

	return nil
}
