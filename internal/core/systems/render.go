package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/entities"
)

type RenderSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
}

func NewRenderSystem(getter componentsGetter, broadcaster Broadcaster) *RenderSystem {
	return &RenderSystem{
		getter:      getter,
		broadcaster: broadcaster,
	}
}

func (s *RenderSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("update")
	for id := range comps {
		var err error

		switch {
		case entities.IsPlayer(id, s.getter):
			err = entities.SendPlayer(id, s.getter, s.broadcaster.SendEntityUpdate)
		case entities.IsShip(id, s.getter):
			err = entities.SendShip(id, s.getter, s.broadcaster.SendEntityUpdate)
		case entities.IsAsteroid(id, s.getter):
			err = entities.SendAsteroid(id, s.getter, s.broadcaster.SendEntityUpdate)
		case entities.IsCosmoAlien(id, s.getter):
			err = entities.SendCosmoAlien(id, s.getter, s.broadcaster.SendEntityUpdate)
		case entities.IsBullet(id, s.getter):
			err = entities.SendBullet(id, s.getter, s.broadcaster.SendEntityUpdate)
		}

		if err != nil {
			return err
		}

	}

	return nil
}
