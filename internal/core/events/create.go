package events

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type CreatePlayerEvent struct {
	ID string
}

func (*CreatePlayerEvent) Type() string { return "createPlayer" }

func NewCreatePlayerEvent(id string) *CreatePlayerEvent {
	return &CreatePlayerEvent{
		ID: id,
	}
}

type CreateShipEvent struct{}

func (*CreateShipEvent) Type() string { return "createShip" }

func NewCreateShipEvent() *CreateShipEvent {
	return &CreateShipEvent{}
}

type CreateAlienEvent struct {
	Position geometry.Vec3
}

func (*CreateAlienEvent) Type() string { return "createAlien" }

func NewCreateAlienEvent(position geometry.Vec3) *CreateAlienEvent {
	return &CreateAlienEvent{
		Position: position,
	}
}

type CreateAsteroidEvent struct {
	Position  geometry.Vec3
	Radius    float64
	Direction geometry.Vec3
	Speed     float64
}

func (*CreateAsteroidEvent) Type() string { return "createAsteroid" }

func NewCreateAsteroidEvent(pos geometry.Vec3, rad float64, dir geometry.Vec3, speed float64) *CreateAsteroidEvent {
	return &CreateAsteroidEvent{
		Position:  pos,
		Radius:    rad,
		Direction: dir,
		Speed:     speed,
	}
}

type CreateCosmoAlienEvent struct {
	Position geometry.Vec3
	ShipID   types.Entity
}

func (*CreateCosmoAlienEvent) Type() string { return "createCosmoAlien" }

func NewCreateCosmoAlienEvent(pos geometry.Vec3, shipID types.Entity) *CreateCosmoAlienEvent {
	return &CreateCosmoAlienEvent{
		Position: pos,
		ShipID:   shipID,
	}
}

type CreateBulletEvent struct {
	Position  geometry.Vec3
	Direction geometry.Vec3
}

func (*CreateBulletEvent) Type() string { return "createBullet" }

func NewCreateBulletEvent(pos, dir geometry.Vec3) *CreateBulletEvent {
	return &CreateBulletEvent{
		Position:  pos,
		Direction: dir,
	}
}

type CreateBreakdownEvent struct {
	Position geometry.Vec3
	WallID   types.Entity
	RoomID   types.Entity
}

func (*CreateBreakdownEvent) Type() string { return "createBreakdown" }

func NewCreateBreakdownEvent(wallID, roomID types.Entity, pos geometry.Vec3) *CreateBreakdownEvent {
	return &CreateBreakdownEvent{
		Position: pos,
		WallID:   wallID,
		RoomID:   roomID,
	}
}

type Door struct {
	ID          string
	RoomA       string
	RoomB       string
	Shape       string
	Center      geometry.Vec3
	HalfExtents geometry.Vec3
	Quaternion  geometry.Quaternion
}

type CreateDoorEvent struct {
	Door Door
}

func (*CreateDoorEvent) Type() string { return "createDoor" }

func NewCreateDoorEvent(door Door) *CreateDoorEvent {
	return &CreateDoorEvent{
		Door: door,
	}
}
