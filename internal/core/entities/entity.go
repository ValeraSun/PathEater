package entities

import (
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	cameraDistant = 2
	cameraHeight  = 1.9
	cameraOffset  = 0.5

	playerSpeed          = 10
	playerMaxHealth      = 100
	playerHalfHeight     = 1
	playerRadius         = 0.7
	playerSpawnX         = 2
	playerSpawnY         = 1
	playerSpawnZ         = 0
	playerAttackCooldown = 1
	playerDamage         = 50

	shipHealth      = 100
	shipSpeed       = 120
	shipBaggage     = 10
	shupWeaponSpeed = 10
	shipAmmo        = 30
	shipDirX        = 1
	shipDirY        = 0

	cosmoAlienRadius = 15
	cosmoAlienSpeed  = 160

	alienHalfHeight       = 1
	alienRadius           = 0.6
	alienAttackDamage     = 0
	alienAttackCooldown   = 2
	alienAttackDistant    = 1
	alienAttackHalfHeight = 1
	alienAttackRadius     = 1
)

type entityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
	AddEntityByID(entity types.Entity, components ...types.Component) error
}

func NewTerminal(adder entityAdder) types.Entity {
	e, _ := adder.AddEntity(
		components.NewInteractableComponent(
			geometry.Vec3{X: 39.85, Y: 1.075, Z: -9},
			geometry.Vec3{X: 0.5, Y: 0.775, Z: 1.150},
			components.InteractTerminal,
		),
	)
	return e
}

func NewPlayer(adder entityAdder, clientID string) types.Entity {
	adder.AddEntityByID(
		types.Entity(clientID),
		components.NewInteractionDetectorComponent(cameraDistant, cameraHeight, cameraOffset),
		components.NewControlShipComponent(),
		components.NewPlayerComponent(),
		components.NewUpdateComponent(),
		components.NewControlComponent(clientID),
		components.NewMovementComponent(playerSpeed, geometry.GetZeroVector()),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewHealthComponent(playerMaxHealth),
		components.NewColliderComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{},
			geometry.Vec3{Y: 1},
			playerHalfHeight,
			playerRadius,
		)),
		components.NewHitboxComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{},
			geometry.Vec3{Y: 1},
			playerHalfHeight,
			playerRadius,
		)),
		components.NewRayAttackComponent(playerAttackCooldown, playerDamage),
		components.NewMovableComponent(),
		components.NewTransformComponent(
			geometry.Vec3{
				X: playerSpawnX,
				Y: playerSpawnY,
				Z: playerSpawnZ,
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
		components.NewHealthComponent(shipHealth),
		components.NewMovementComponent(shipSpeed, geometry.Vec3{X: shipDirX, Y: shipDirY, Z: 0}),
		components.NewVelocityComponent(),
		components.NewShipComponent(shipBaggage),
		components.NewWeaponComponent(geometry.GetZeroVector(), shupWeaponSpeed, shipAmmo),
		components.NewColliderComponent(geometry.NewTriangleCollider(
			geometry.Vec3{X: 20, Y: 0, Z: 0},
			geometry.Vec3{X: -20, Y: 20, Z: 0},
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
		components.NewUpdateComponent(),
		components.NewNavigationEntityComponent(),
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

func NewCosmoAlien(adder entityAdder, pos geometry.Vec3) types.Entity {
	e, _ := adder.AddEntity(
		components.NewTransformComponent(
			pos,
			geometry.GetZeroVector(),
		),
		components.NewMovementComponent(cosmoAlienSpeed, geometry.GetZeroVector()),
		components.NewExternalVelocityComponent(),
		components.NewVelocityComponent(),
		components.NewUpdateComponent(),
		components.NewNavigationEntityComponent(),
		components.NewCosmoAlienComponent(),
		components.NewHealthComponent(10),
		components.NewColliderComponent(geometry.NewCircleCollider(
			pos,
			cosmoAlienRadius,
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
		components.NewNavigationEntityComponent(),
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

func NewAlien(adder entityAdder, position geometry.Vec3) types.Entity {
	alien, _ := adder.AddEntity(
		components.NewTransformComponent(position, geometry.Vec3{}),
		components.NewVisionComponent(),
		components.NewAttackComponent(
			alienAttackDamage,
			alienAttackCooldown,
			alienAttackDistant,
			geometry.NewCapsuleCollider(
				geometry.GetZeroVector(),
				geometry.Vec3{Y: 1},
				alienAttackHalfHeight,
				alienAttackRadius,
			),
		),
		components.NewMovementComponent(2, geometry.GetZeroVector()),
		components.NewMovableComponent(),
		components.NewVelocityComponent(),
		components.NewUpdateComponent(),
		components.NewHealthComponent(100),
		components.NewAIComponent(),
		components.NewColliderComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{},
			geometry.Vec3{Y: 2},
			alienHalfHeight,
			alienRadius,
		)),
		components.NewHitboxComponent(geometry.NewCapsuleCollider(
			geometry.Vec3{},
			geometry.Vec3{Y: 2},
			alienHalfHeight,
			alienRadius,
		)),
	)
	return alien
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
	Health   int           `json:"health"`
}

type shipData struct {
	Collider        geometry.Collider `json:"collider"`
	WeaponDirection geometry.Vec3     `json:"weapon_direction"`
	ShootSuccess    bool              `json:"shoot_success"`
	BaggageStatus   int               `json:"baggage_status"`
	Health          int               `json:"health"`
}

type alienData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
	Health   int           `json:"health"`
	Dead     bool          `json:"dead"`
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
	return getter.HasComponents(id, "transform", "control")
}

func SendPlayer(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	c, _ := getter.GetComponent(id, "transform")
	transform := c.(*components.TransformComponent)

	c, _ = getter.GetComponent(id, "control")
	control := c.(*components.ControlComponent)

	c, _ = getter.GetComponent(id, "health")
	hp := c.(*components.HealthComponent)

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "player",
		Data: playerData{
			Position: transform.Position,
			Rotation: control.Direction,
			Health:   hp.Health,
		},
	},
	)

	return nil
}

func SendAlien(id types.Entity, getter componentsGetter, broadcaster Sendler) error {

	c, _ := getter.GetComponent(id, "transform")
	transform := c.(*components.TransformComponent)

	c, _ = getter.GetComponent(id, "ai")
	ai := c.(*components.AIComponent)

	c, _ = getter.GetComponent(id, "health")
	hp := c.(*components.HealthComponent)

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "alien",
		Data: alienData{
			Position: transform.Position,
			Rotation: ai.Direction,
			Health:   hp.Health,
			Dead:     hp.Health == 0,
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
		return errors.New("Не найден необходимый копнонент")
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

func SendAsteroid(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	/*shipPos, err := getShipPosition(getter)
	if err != nil {
		return err
	}*/

	c, ok := getter.GetComponent(id, "asteroid")
	if !ok {
		return errors.New("не найден компонент asteroid")
	}
	aster := c.(*components.AsteroidComponent)

	c, ok = getter.GetComponent(id, "transform")
	if !ok {
		return errors.New("не найден компонент transform")
	}
	transform := c.(*components.TransformComponent)

	c, ok = getter.GetComponent(id, "collider")
	if !ok {
		return errors.New("не найден компонент collider")
	}
	col := c.(*components.ColliderComponent)

	/*relPos := geometry.Vec3{
		X: transform.Position.X - shipPos.X,
		Y: transform.Position.Y - shipPos.Y,
	}*/

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "asteroid",
		Data: asteroidData{
			Position:  transform.Position,
			Radius:    col.Collider.(*geometry.CircleCollider).Radius,
			Destroyed: aster.Destroyed,
		},
	})

	return nil
}

func IsCosmoAlien(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "cosmoAlien")
}

func SendCosmoAlien(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	c, ok := getter.GetComponent(id, "cosmoAlien")
	if !ok {
		return errors.New("не найден компонент cosmoAlien")
	}
	alien := c.(*components.CosmoAlienComponent)

	c, ok = getter.GetComponent(id, "transform")
	if !ok {
		return errors.New("не найден компонент transform")
	}
	transform := c.(*components.TransformComponent)

	c, ok = getter.GetComponent(id, "health")
	if !ok {
		return errors.New("не найден компонент health")
	}
	hp := c.(*components.HealthComponent)

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "cosmoAlien",
		Data: cosmoAlienData{
			Position: transform.Position,
			Rotation: transform.Direction,
			Health:   hp.Health,
			Died:     alien.Died,
		},
	})

	return nil
}

func IsBullet(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "bullet")
}

func SendBullet(id types.Entity, getter componentsGetter, broadcaster Sendler) error {
	shipPos, err := getShipPosition(getter)
	if err != nil {
		return err
	}

	c, ok := getter.GetComponent(id, "transform")
	if !ok {
		return errors.New("не найден компонент transform")
	}
	transform := c.(*components.TransformComponent)

	relPos := geometry.Vec3{
		X: transform.Position.X - shipPos.X,
		Y: transform.Position.Y - shipPos.Y,
		Z: transform.Position.Z - shipPos.Z,
	}

	broadcaster(ecs.EntityInfo{
		ID:   id,
		Type: "bullet",
		Data: bulletData{
			Position: relPos,
		},
	})

	return nil
}

func getShipPosition(getter componentsGetter) (geometry.Vec3, error) {
	entities := getter.GetEntitiesByComponent("ship")
	for id := range entities {
		if comp, ok := getter.GetComponent(id, "transform"); ok {
			trans := comp.(*components.TransformComponent)
			return trans.Position, nil
		}
	}
	return geometry.Vec3{}, errors.New("ship not found")
}
