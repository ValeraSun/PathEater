package sendler

import (
	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}

type EntityInfo struct {
	ID   types.Entity `json:"id"`
	Type string       `json:"type"`
	Data any          `json:"data"`
}

type sendler struct {
	componentsGetter
}

func New(componentsChecker componentsGetter) *sendler {
	return &sendler{
		componentsChecker,
	}
}

func (s *sendler) Send(id types.Entity, broadcaster func(EntityInfo) error) error {
	switch {
	case entities.IsPlayer(id, s):
		s.sendPlayer(id, broadcaster)

	case entities.IsAsteroid(id, s):
		s.sendAsteroid(id, broadcaster)

	case entities.IsShip(id, s):
		s.sendShip(id, broadcaster)

	case entities.IsCosmoAlien(id, s):
		s.sendCosmoAlien(id, broadcaster)

	case entities.IsAlien(id, s):
		s.sendAlien(id, broadcaster)

	case entities.IsBullet(id, s):
		s.sendBullet(id, broadcaster)

	case entities.IsDoor(id, s):
		s.sendDoor(id, broadcaster)
	}
	return nil
}

type GameOverInfo struct {
	Win    bool `json:"win"`
	Status int  `json:"status"`
}

type TimeInfo struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
