package events

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type DeletePlayerEvent struct {
	ID types.Entity
}

func (*DeletePlayerEvent) Type() string { return "deletePlayer" }

func NewDeletePlayerEvent(id types.Entity) *DeletePlayerEvent {
	return &DeletePlayerEvent{
		ID: id,
	}
}

type DeleteShipEvent struct {
	ID types.Entity
}

func (*DeleteShipEvent) Type() string { return "deleteShip" }

func NewDeleteShipEvent(id types.Entity) *DeleteShipEvent {
	return &DeleteShipEvent{
		ID: id,
	}
}

type DeleteAsteroidEvent struct {
	ID types.Entity
}

func (*DeleteAsteroidEvent) Type() string { return "deleteAsteroid" }

func NewDeleteAsteroidEvent(id types.Entity) *DeleteAsteroidEvent {
	return &DeleteAsteroidEvent{
		ID: id,
	}
}

type DeleteCosmoAlienEvent struct {
	ID types.Entity
}

func (*DeleteCosmoAlienEvent) Type() string { return "deleteCosmoAlien" }

func NewDeleteCosmoAlienEvent(id types.Entity) *DeleteCosmoAlienEvent {
	return &DeleteCosmoAlienEvent{
		ID: id,
	}
}

type DeleteBulletEvent struct {
	ID types.Entity
}

func (*DeleteBulletEvent) Type() string { return "deleteBullet" }

func NewDeleteBulletEvent(id types.Entity) *DeleteBulletEvent {
	return &DeleteBulletEvent{
		ID: id,
	}
}
